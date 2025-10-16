<script lang="ts">
    import { onDestroy } from "svelte";
    import { getUserToken } from "../../lib/states/auth_state.svelte";
    import { WS_URL } from "../../lib/constants";

    let logContainer: HTMLDivElement;
    let ws: WebSocket | null = $state(null);
    let logMessages = $state<string[]>([]);
    let connectionStatus = $state<
        "disconnected" | "connecting" | "connected" | "error"
    >("disconnected");

    function log(message: string) {
        logMessages = [
            ...logMessages,
            `[${new Date().toLocaleTimeString()}] ${message}`,
        ];
    }

    $effect(() => {
        if (logContainer && logMessages.length) {
            logContainer.scrollTop = logContainer.scrollHeight;
        }
    });

    function connect() {
        if (ws?.readyState === WebSocket.OPEN) {
            log("Already connected.");
            return;
        }

        connectionStatus = "connecting";
        const token = getUserToken();

        if (!token) {
            log("Error: No authentication token found");
            connectionStatus = "error";
            return;
        }

        try {
            ws = new WebSocket(`${WS_URL}/admin/ws/updater`);

            ws.onopen = () => {
                connectionStatus = "connected";
                log("Connected to /ws/updater");
            };

            ws.onmessage = (event) => {
                log(`[Update] ${event.data}`);
            };

            ws.onclose = () => {
                connectionStatus = "disconnected";
                log("Connection closed");
            };

            ws.onerror = (error) => {
                connectionStatus = "error";
                log(`Error: ${error.type}`);
            };
        } catch (error) {
            connectionStatus = "error";
            log(`Connection failed: ${error}`);
        }
    }

    function startUpdater() {
        if (ws?.readyState === WebSocket.OPEN) {
            ws.send("START");
            log("> START command sent");
        } else {
            log("WebSocket not connected.");
        }
    }

    function stopUpdater() {
        if (ws?.readyState === WebSocket.OPEN) {
            ws.send("STOP");
            log("> STOP command sent");
        } else {
            log("WebSocket not connected.");
        }
    }

    function disconnect() {
        if (ws) {
            ws.close();
            ws = null;
            connectionStatus = "disconnected";
            log("Disconnected manually");
        }
    }

    function clearLog() {
        logMessages = [];
    }

    onDestroy(() => {
        if (ws) {
            ws.close();
        }
    });
</script>

<div class="p-4 sm:p-6">
    <!-- Header -->
    <div
        class="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4 mb-6"
    >
        <h2 class="text-xl font-bold text-[#19183B] text-center sm:text-left">
            Updater Control Panel
        </h2>
        <div class="flex flex-wrap gap-2 justify-center sm:justify-end">
            <button
                class="flex items-center justify-center gap-2 cursor-pointer px-4 py-2 bg-[#19183B] text-white rounded-lg hover:bg-[#2a2852] transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                onclick={connect}
                disabled={connectionStatus === "connecting" ||
                    connectionStatus === "connected"}
            >
                {#if connectionStatus === "connecting"}
                    <!-- Loading Spinner -->
                    <svg
                        class="animate-spin h-5 w-5 text-white"
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                    >
                        <circle
                            class="opacity-25"
                            cx="12"
                            cy="12"
                            r="10"
                            stroke="currentColor"
                            stroke-width="4"
                        ></circle>
                        <path
                            class="opacity-75"
                            fill="currentColor"
                            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                        ></path>
                    </svg>
                    <span class="sm:hidden">Connecting...</span>
                {:else}
                    <svg
                        class="w-5 h-5"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M5 12h14M12 5l7 7-7 7"
                        ></path>
                    </svg>
                    <span class="sm:hidden">Connect</span>
                {/if}
            </button>

            <button
                class="flex items-center justify-center gap-2 cursor-pointer px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                onclick={startUpdater}
                disabled={connectionStatus !== "connected"}
            >
                <svg
                    class="w-5 h-5"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"
                    ></path>
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                    ></path>
                </svg>
                <span class="sm:hidden">Start</span>
            </button>

            <button
                class="flex items-center justify-center gap-2 cursor-pointer px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                onclick={stopUpdater}
                disabled={connectionStatus !== "connected"}
            >
                <svg
                    class="w-5 h-5"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                    ></path>
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M9 10a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1h-4a1 1 0 01-1-1v-4z"
                    ></path>
                </svg>
                <span class="sm:hidden">Stop</span>
            </button>

            <button
                class="flex items-center justify-center gap-2 cursor-pointer px-4 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                onclick={disconnect}
                disabled={connectionStatus !== "connected"}
            >
                <svg
                    class="w-5 h-5"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
                    ></path>
                </svg>
                <span class="sm:hidden">Disconnect</span>
            </button>

            <button
                class="flex items-center justify-center gap-2 cursor-pointer px-4 py-2 bg-[#19183B] text-white rounded-lg hover:bg-[#2a2852] transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                onclick={clearLog}
            >
                <svg
                    class="w-5 h-5"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                    ></path>
                </svg>
                <span class="sm:hidden">Clear Log</span>
            </button>
        </div>
    </div>

    <!-- Connection Status Display -->
    <div class="mb-6 p-4 rounded-lg bg-gray-50 border">
        <div class="flex items-center justify-between">
            <div class="flex items-center space-x-3">
                <div
                    class={`w-3 h-3 rounded-full ${
                        connectionStatus === "connected"
                            ? "bg-green-500"
                            : connectionStatus === "connecting"
                              ? "bg-yellow-500"
                              : connectionStatus === "error"
                                ? "bg-red-500"
                                : "bg-gray-500"
                    }`}
                ></div>
                <span class="font-medium text-[#19183B]">
                    Status:
                    {connectionStatus === "connected"
                        ? "Connected"
                        : connectionStatus === "connecting"
                          ? "Connecting..."
                          : connectionStatus === "error"
                            ? "Connection Error"
                            : "Disconnected"}
                </span>
            </div>
        </div>
    </div>

    <!-- Mobile Card View -->
    <div class="block sm:hidden space-y-4">
        <!-- Control Info Card -->
        <div class="bg-white rounded-lg border border-gray-200 p-4 shadow-sm">
            <div class="text-center">
                <!-- Quick Status -->
                <div class="bg-gray-50 rounded-lg p-3">
                    <div class="flex justify-between items-center text-sm">
                        <span class="text-[#708993]">Connection:</span>
                        <span
                            class={`font-medium ${
                                connectionStatus === "connected"
                                    ? "text-green-600"
                                    : connectionStatus === "connecting"
                                      ? "text-yellow-600"
                                      : connectionStatus === "error"
                                        ? "text-red-600"
                                        : "text-gray-600"
                            }`}
                        >
                            {connectionStatus === "connected"
                                ? "Active"
                                : connectionStatus === "connecting"
                                  ? "Connecting"
                                  : connectionStatus === "error"
                                    ? "Error"
                                    : "Inactive"}
                        </span>
                    </div>
                </div>
            </div>
        </div>

        <!-- Log Display Card -->
        <div class="bg-white rounded-lg border border-gray-200 p-4 shadow-sm">
            <h3 class="text-lg font-semibold text-[#19183B] mb-3">
                Update Log
            </h3>
            <div
                class="h-64 bg-[#1e1e1e] text-green-400 font-mono text-sm p-3 rounded border overflow-y-auto whitespace-pre-wrap"
            >
                {#if logMessages.length === 0}
                    <div class="text-gray-500 text-center py-8">
                        <svg
                            class="mx-auto h-8 w-8 text-gray-400 mb-2"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                            ></path>
                        </svg>
                        <div>Waiting for connection...</div>
                    </div>
                {:else}
                    {#each logMessages as message}
                        <div
                            class="py-1 border-b border-gray-700 last:border-b-0"
                        >
                            {message}
                        </div>
                    {/each}
                {/if}
            </div>
        </div>
    </div>

    <!-- Desktop Table View -->
    <div class="hidden sm:block">
        <div
            class="bg-white rounded-lg border border-gray-200 shadow-sm overflow-hidden"
        >
            <!-- Table Header -->
            <div class="bg-gray-50 px-6 py-4 border-b border-gray-200">
                <div class="flex justify-between items-center">
                    <h3 class="text-lg font-semibold text-[#19183B]">
                        Real-time Update Log
                    </h3>
                    <div class="text-sm text-[#708993]">
                        {logMessages.length} log entries
                    </div>
                </div>
            </div>

            <!-- Log Content -->
            <div class="p-0">
                <div
                    bind:this={logContainer}
                    class="h-96 bg-[#1e1e1e] text-green-400 font-mono text-sm p-6 overflow-y-auto whitespace-pre-wrap"
                >
                    {#if logMessages.length === 0}
                        <div
                            class="h-full text-gray-500 text-center flex flex-col justify-center items-center"
                        >
                            <svg
                                class="mx-auto h-12 w-12 text-gray-400 mb-4"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                                ></path>
                            </svg>
                            <p class="text-gray-500">
                                Connect to the updater service to see real-time <span
                                    >logs</span
                                >
                            </p>
                        </div>
                    {:else}
                        {#each logMessages as message}
                            <div
                                class="py-2 border-b border-gray-700 last:border-b-0"
                            >
                                {message}
                            </div>
                        {/each}
                    {/if}
                </div>
            </div>
        </div>
    </div>
</div>
