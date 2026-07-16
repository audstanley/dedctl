package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"dedctl/internal/config"
	"dedctl/internal/service"
)

func NewCacheImagesCmd() *cobra.Command {
	var force bool

	cacheCmd := &cobra.Command{
		Use:   "cache-images [game]...",
		Short: "Cache game cover images from SteamDB",
		Long: `Downloads game cover images from SteamDB for the specified games.

If no game names are provided, caches images for all games in the metadata.
Games must have an app_id set in the metadata.yaml file.

Examples:
  dedctl cache-images                # cache all games in metadata
  dedctl cache-images csgo rust      # cache specific games
  dedctl cache-images --force        # re-download all images even if cached`,
		Run: func(cmd *cobra.Command, args []string) {
			metaDir := resolveMetadataDir()
			imgDir := filepath.Join(metaDir, "img")

			// Ensure metadata.yaml exists
			metaPath := filepath.Join(metaDir, "metadata.yaml")
			if _, err := os.Stat(metaPath); os.IsNotExist(err) {
				fmt.Printf("Creating metadata.yaml at %s\n", metaPath)
				if err := config.SaveMetadata(metaDir, &config.Metadata{Games: make(map[string]config.GameMetadata)}); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(1)
				}
			}

			// Load metadata
			meta, err := config.LoadMetadata(metaDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading metadata: %v\n", err)
				os.Exit(1)
			}

			imageService := service.NewImageService()

			// Determine which games to cache
			var gamesToCache []string
			if len(args) > 0 {
				gamesToCache = args
			} else {
				gamesToCache = meta.SortedGames()
			}

			// Filter to games with app_id
			var withAppId []string
			for _, name := range gamesToCache {
				appId := meta.GetAppId(name)
				if appId <= 0 {
					fmt.Printf("Skipping '%s': no app_id set\n", name)
					continue
				}
				withAppId = append(withAppId, name)
			}

			if len(withAppId) == 0 {
				fmt.Println("No games with app_id found. Set app_id in metadata.yaml first.")
				os.Exit(0)
			}

			// Build metadata map for the image service
			metaMap := make(map[string]struct{ AppId int })
			for name, gm := range meta.Games {
				metaMap[name] = struct{ AppId int }{AppId: gm.AppId}
			}

			// Ensure img directory exists
			if err := os.MkdirAll(imgDir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error: failed to create image directory: %v\n", err)
				os.Exit(1)
			}

			// Download images
			fmt.Printf("Caching images for %d game(s)...\n", len(withAppId))
			success := 0
			failed := 0

			for _, name := range withAppId {
				appId := meta.GetAppId(name)
				if force || !imageService.ImageExists(appId, imgDir) {
					if err := imageService.CacheSingleImage(name, appId, imgDir); err != nil {
						fmt.Printf("  FAILED %s (app %d): %v\n", name, appId, err)
						failed++
					} else {
						fmt.Printf("  OK     %s (app %d)\n", name, appId)
						success++
					}
				} else {
					fmt.Printf("  SKIP   %s (app %d) - already cached\n", name, appId)
				}
			}

			fmt.Printf("\nDone! %d cached, %d failed\n", success, failed)
			if failed > 0 {
				os.Exit(1)
			}
		},
	}

	cacheCmd.Flags().BoolVar(&force, "force", false, "force re-download of cached images")
	return cacheCmd
}
