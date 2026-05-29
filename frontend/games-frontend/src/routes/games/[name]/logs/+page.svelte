<script lang="ts">
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { onMount, onDestroy } from 'svelte';
  import { api } from '$lib/api/client';

  let { name } = $props();
  let logs = $state<string[]>([]);
  let autoScroll = $state(true);
  let eventSource: EventSource | null = null;
  let connected = $state(false);
  let reconnecting = $state(false);
  let reconnectCount = $state(0);
  let maxRetries = $state(5);

  $effect(() => {
    if (autoScroll) {
      scrollContainer?.scrollTo(0, scrollContainer.scrollHeight);
    }
  });

  let scrollContainer: HTMLDivElement | null = null;

  function connect() {
    if (eventSource) {
      eventSource.close();
    }

    const source = api.streamLogs(name);
    eventSource = source;
    reconnecting = false;

    source.onopen = () => {
      connected = true;
      reconnectCount = 0;
    };

    source.onmessage = (event) => {
      logs = [...logs, event.data];
    };

    source.onerror = (error) => {
      console.error('SSE Error:', error);
      connected = false;
      source.close();
      eventSource = null;

      if (reconnectCount < maxRetries) {
        reconnectCount++;
        reconnecting = true;
        setTimeout(connect, 2000 * reconnectCount);
      }
    };
  }

  onMount(() => {
    connect();
  });

  onDestroy(() => {
    eventSource?.close();
  });

  function handleScroll(e: Event) {
    const target = e.target as HTMLDivElement;
    const isNearBottom = target.scrollHeight - target.scrollTop - target.clientHeight < 100;
    autoScroll = isNearBottom;
  }

  function handleRefresh() {
    if (eventSource !== null) {
      eventSource.close();
    }
    logs = [];
    reconnectCount = 0;
    connect();
  }

  function handleClear() {
    logs = [];
  }

  function handleGoBack() {
    goto(`/games/${name}`);
  }

  function formatTimestamp(timestamp: string): string {
    try {
      const date = new Date(parseInt(timestamp));
      return date.toLocaleTimeString();
    } catch {
      return timestamp;
    }
  }

  function extractMessage(line: string): string {
    const match = line.match(/\] (.+)$/);
    return match ? match[1] : line;
  }
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <div class="flex items-center space-x-4">
      <button
        onclick={handleGoBack}
        class="bg-gray-700 hover:bg-gray-600 text-white px-4 py-2 rounded transition"
      >
        ← Back to Game
      </button>
      <h1 class="text-3xl font-bold text-white capitalize">Logs - {name}</h1>
    </div>

    <div class="flex items-center space-x-3">
      {#if connected}
        <span class="text-green-400 text-sm flex items-center">
          <span class="w-2 h-2 bg-green-400 rounded-full mr-2 animate-pulse"></span>
          Connected
        </span>
      {:else if reconnecting}
        <span class="text-yellow-400 text-sm">
          Reconnecting... ({reconnectCount}/{maxRetries})
        </span>
      {:else}
        <span class="text-red-400 text-sm">Disconnected</span>
      {/if}

      <button
        onclick={handleRefresh}
        class="bg-blue-600 hover:bg-blue-700 text-white px-3 py-2 rounded transition"
        title="Refresh logs"
      >
        ⟳
      </button>

      <button
        onclick={handleClear}
        class="bg-gray-700 hover:bg-gray-600 text-white px-3 py-2 rounded transition"
        title="Clear logs"
      >
        Clear
      </button>
    </div>
  </div>

  <div class="bg-gray-800 rounded-lg shadow-lg overflow-hidden">
    <div class="bg-gray-750 px-4 py-2 border-b border-gray-700 flex items-center justify-between">
      <span class="text-gray-400 text-sm">Real-time stream</span>
      <label class="flex items-center space-x-2 text-gray-400 text-sm cursor-pointer">
        <input
          type="checkbox"
          bind:checked={autoScroll}
          class="rounded bg-gray-700 border-gray-600 text-blue-600"
        />
        <span>Auto-scroll</span>
      </label>
    </div>

    <div
      bind:this={scrollContainer}
      onscroll={handleScroll}
      class="h-96 overflow-y-auto p-4 font-mono text-sm space-y-1"
    >
      {#if logs.length === 0}
        <div class="text-gray-500 text-center py-8">
          <p>No logs yet. Waiting for server output...</p>
        </div>
      {:else}
        {#each logs as log (log)}
          <div class="text-gray-300 hover:bg-gray-750 px-2 py-0.5 rounded">
            <span class="text-blue-400 mr-2">
              {#if log.includes('[') && log.includes(']')}
                {@html formatTimestamp(extractMessage(log))}
              {:else}
                -
              {/if}
            </span>
            {@html extractMessage(log)}
          </div>
        {/each}
      {/if}
    </div>
  </div>
</div>
