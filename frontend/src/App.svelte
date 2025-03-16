
<script lang="ts">
  let pricing = [
    { module: 'Cuper', cost: '0.4'},
    { module: 'Engel', cost: '0.6'}
  ];
  let coordsUrl = "https://olimp.miet.ru/ppo_it/api";
  let listenerCoords: number[] = [];
  let senderCoords: number[] = [];
  
  // Add this function to fetch coordinates and pricing data
  async function fetchCoordsData() {
    try {
      const response = await fetch("http://localhost:8080/api/coords", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({ url: coordsUrl })
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch coordinates data: ${response.status}`);
      }

      const data = await response.json();
      
      if (data) {
        // Update listener and sender coordinates
        listenerCoords = data.listener || [];
        senderCoords = data.sender || [];
        
        // Update pricing information
        if (data.price) {
          pricing = [
            { module: 'Cuper', cost: data.price.cuper.toString() },
            { module: 'Engel', cost: data.price.engel.toString() }
          ];
        }
        
        // Redraw the grid to show updated information
        renderMarsGrid();
      }
    } catch (error) {
      console.error("Error fetching coordinates data:", error);
    }
  }

  let canvas: HTMLCanvasElement;
  
  const modules = [
    { id: 1, name: 'Cuper', status: 'Not Used', quantyt: '0' },
    { id: 2, name: 'Engel', status: 'Not Used', quantyty: '0' }
  ];

  
  
  
  import { onMount } from 'svelte';
  
  
  let marsGrid: number[][] = [];
  let chunks: { x: number, y: number, data: number[][] }[] = [];
  let selectedChunk: number | null = null;
  let targetChunk: number | null = null;
  
  
  function initMarsGrid() {
    
    marsGrid = Array(256).fill(0).map(() => Array(256).fill(0));
    
    
    chunks = [];
    
    for (let chunkY = 0; chunkY < 4; chunkY++) {
      for (let chunkX = 0; chunkX < 4; chunkX++) {
        const chunkData: number[][] = [];
        
        for (let y = 0; y < 64; y++) {
          chunkData.push([]);
          for (let x = 0; x < 64; x++) {
            //, will be replaced with API data
            chunkData[y].push(0);
          }
        }

        //, will be updated with API data
        chunks.push({
          x: chunkX,
          y: chunkY,
          data: chunkData,
        });
      }
    }

    //, fetch real data from API
    fetchChunkDataFromAPI();
  }

  //
  const chunkUrls = [
    "https://olimp.miet.ru/ppo_it/api",
    "https://olimp.miet.ru/ppo_it/api",
    "https://olimp.miet.ru/ppo_it/api",
    "https://olimp.miet.ru/ppo_it/api",
    "https://olimp.miet.ru/ppo_it/api",
    "https://olimp.miet.ru/ppo_it/api",
    "https://olimp.miet.ru/ppo_it/api",
    "https://olimp.miet.ru/ppo_it/api",
    "https://olimp.miet.ru/ppo_it/api",
    "https://olimp.miet.ru/ppo_it/api",
    "https://olimp.miet.ru/ppo_it/api",
    "https://olimp.miet.ru/ppo_it/api",
    "https://olimp.miet.ru/ppo_it/api",
    "https://olimp.miet.ru/ppo_it/api",
    "https://olimp.miet.ru/ppo_it/api",
    "https://olimp.miet.ru/ppo_it/api"
  ];

    // Function to update a specific chunk URL
    function updateChunkUrl(index: number, newUrl: string) {
    if (index >= 0 && index < chunkUrls.length) {
      chunkUrls[index] = newUrl;
    }
  }

  // Function to refresh data for a specific chunk
  async function refreshChunkData(index: number) {
    if (index >= 0 && index < chunks.length) {
      try {
        const response = await fetch("http://localhost:8080/api/map", {
          method: "POST",
          headers: {
            "Content-Type": "application/json"
          },
          body: JSON.stringify({ url: chunkUrls[index] })
        });

        if (!response.ok) {
          throw new Error(`Failed to fetch chunk data: ${response.status}`);
        }

        const data = await response.json();

        if (data && data.map) {
          const decodedData: number[][] = [];
          
          for (let i = 0; i < data.map.length; i++) {
            const line = data.map[i];
            const decodedLine = decodeBase64Line(line);
            decodedData.push(decodedLine);
          }
          
          chunks[index].data = decodedData;
          updateMarsGridFromChunks();
        }
      } catch (error) {
        console.error(`Error refreshing chunk ${index}:`, error);
      }
    }
  }

  // Function to refresh all chunks data
  function refreshAllChunks() {
    fetchChunkDataFromAPI();
  }

  async function fetchChunkDataFromAPI() {
    try {
      //
      const fetchPromises = chunks.map(async (chunk, index) => {
        const response = await fetch("http://localhost:8080/api/map", {
          method: "POST",
          headers: {
            "Content-Type": "application/json"
          },
          body: JSON.stringify({ url: chunkUrls[index] })
        });

        if (!response.ok) {
          throw new Error(`Failed to fetch chunk data: ${response.status}`);
        }

        const data = await response.json();

        //
                  //
          //
          if (data && data.map) {
            function decodeBase64Line(base64String: string): number[] {
    //
    const binaryString = atob(base64String);
    
    //
    const result = new Uint8Array(binaryString.length);
    for (let i = 0; i < binaryString.length; i++) {
      result[i] = binaryString.charCodeAt(i);
    }
    
    //
    return Array.from(result);
  }
            const decodedData: number[][] = [];
            
            //
            for (let i = 0; i < data.map.length; i++) {
              const line = data.map[i];
              //
              const decodedLine = decodeBase64Line(line);
              decodedData.push(decodedLine);
            }
            
            chunks[index].data = decodedData;
          }
      });

      //
      await Promise.all(fetchPromises);

      //
      updateMarsGridFromChunks();
    } catch (error) {
      console.error("Error fetching chunk data:", error);
      //
      chunks.forEach(chunk => {
        for (let y = 0; y < 64; y++) {
          for (let x = 0; x < 64; x++) {
            chunk.data[y][x] = Math.floor(Math.random() * 256);
          }
        }
      });
      updateMarsGridFromChunks();
    }
  }
  
  
  function updateMarsGridFromChunks() {
    
    marsGrid = Array(256).fill(0).map(() => Array(256).fill(0));
    
    
    chunks.forEach(chunk => {
      const startX = chunk.x * 64;
      const startY = chunk.y * 64;
      
      for (let y = 0; y < 64; y++) {
        for (let x = 0; x < 64; x++) {
          if (startX + x < 256 && startY + y < 256) {
            marsGrid[startY + y][startX + x] = chunk.data[y][x];
          }
        }
      }
    });
    
    renderMarsGrid();
  }
  
  
  function swapChunks(index1: number, index2: number) {
    if (index1 === index2) return;
    
    
    const tempX = chunks[index1].x;
    const tempY = chunks[index1].y;
    
    chunks[index1].x = chunks[index2].x;
    chunks[index1].y = chunks[index2].y;
    
    chunks[index2].x = tempX;
    chunks[index2].y = tempY;
    
    updateMarsGridFromChunks();
    selectedChunk = null;
    targetChunk = null;
  }
  
  
  function handleChunkClick(event: MouseEvent) {
    if (!canvas) return;
    
    const rect = canvas.getBoundingClientRect();
    const x = event.clientX - rect.left;
    const y = event.clientY - rect.top;
    
    
    const gridX = Math.floor(x / (canvas.width / 256));
    const gridY = Math.floor(y / (canvas.height / 256));
    
    
    const chunkX = Math.floor(gridX / 64);
    const chunkY = Math.floor(gridY / 64);
    
    const clickedChunkIndex = chunks.findIndex(
      chunk => chunk.x === chunkX && chunk.y === chunkY
    );
    
    if (clickedChunkIndex !== -1) {
      if (selectedChunk === null) {
        
        selectedChunk = clickedChunkIndex;
      } else {
        
        targetChunk = clickedChunkIndex;
        swapChunks(selectedChunk, targetChunk);
      }
      
      renderMarsGrid();
    }
  }
  
  
  function renderMarsGrid() {
  if (!canvas) return;
  
  const ctx = canvas.getContext('2d');
  if (!ctx) return;
  
  ctx.clearRect(0, 0, canvas.width, canvas.height);
  
  const cellWidth = canvas.width / 256;
  const cellHeight = canvas.height / 256;
  
  //
  for (let y = 0; y < 256; y++) {
    for (let x = 0; x < 256; x++) {
      const elevation = marsGrid[y][x];
      
      const intensity = Math.floor(50 + (elevation / 255) * 205); 
      ctx.fillStyle = `rgb(${intensity}, ${Math.floor(intensity/3)}, ${Math.floor(intensity/3)})`;
      ctx.fillRect(x * cellWidth, y * cellHeight, cellWidth, cellHeight);
    }
  }
  
  //
  ctx.strokeStyle = '#ffffff';
  ctx.lineWidth = 2;
  
  //
  for (let i = 1; i < 4; i++) {
    const y = i * 64 * cellHeight;
    ctx.beginPath();
    ctx.moveTo(0, y);
    ctx.lineTo(canvas.width, y);
    ctx.stroke();
  }
  
  //
  for (let i = 1; i < 4; i++) {
    const x = i * 64 * cellWidth;
    ctx.beginPath();
    ctx.moveTo(x, 0);
    ctx.lineTo(x, canvas.height);
    ctx.stroke();
  }
  
  //
  if (selectedChunk !== null) {
    const chunk = chunks[selectedChunk];
    ctx.strokeStyle = '#ff5733';
    ctx.lineWidth = 3;
    ctx.strokeRect(
      chunk.x * 64 * cellWidth,
      chunk.y * 64 * cellHeight,
      64 * cellWidth,
      64 * cellHeight
    );

    // Add this code at the end of the renderMarsGrid function, just before the final closing brace
// Mark listener and sender positions if available
if (listenerCoords.length === 2) {
  const [x, y] = listenerCoords;
  ctx.fillStyle = 'blue';
  ctx.beginPath();
  ctx.arc(
    x * cellWidth + cellWidth/2, 
    y * cellHeight + cellHeight/2, 
    Math.max(cellWidth, cellHeight) * 3, 
    0, 
    2 * Math.PI
  );
  ctx.fill();
  ctx.fillStyle = 'white';
  ctx.font = '12px Arial';
  ctx.fillText('L', x * cellWidth + cellWidth/2 - 4, y * cellHeight + cellHeight/2 + 4);
}

if (senderCoords.length === 2) {
  const [x, y] = senderCoords;
  ctx.fillStyle = 'green';
  ctx.beginPath();
  ctx.arc(
    x * cellWidth + cellWidth/2, 
    y * cellHeight + cellHeight/2, 
    Math.max(cellWidth, cellHeight) * 3, 
    0, 
    2 * Math.PI
  );
  ctx.fill();
  ctx.fillStyle = 'white';
  ctx.font = '12px Arial';
  ctx.fillText('S', x * cellWidth + cellWidth/2 - 4, y * cellHeight + cellHeight/2 + 4);
}
  }
  
  //
  ctx.strokeStyle = '#888';
  ctx.lineWidth = 2;
  ctx.strokeRect(0, 0, canvas.width, canvas.height);
}
  
  onMount(() => {
    if (canvas) {
      canvas.width = 512;
      canvas.height = 512;
      
      
      initMarsGrid();
      
      
      canvas.addEventListener('click', handleChunkClick);
      
      const handleResize = () => {
        canvas.width = canvas.parentElement!.clientWidth;
        canvas.height = canvas.parentElement!.clientHeight;
        renderMarsGrid();
      };
      
      window.addEventListener('resize', handleResize);
      
      return () => {
        window.removeEventListener('resize', handleResize);
        canvas.removeEventListener('click', handleChunkClick);
      };
    }
  });
  
  
  function resetGrid() {
    initMarsGrid();
  }
</script>

<main class="min-h-screen bg-base-200 p-4">
  <div class="container mx-auto">
    <h1 class="text-3xl font-bold mb-6 text-center">Mars Surface Module Placement System</h1>
    
    <div class="flex flex-col lg:flex-row gap-4 responsive-layout">
      <div class="flex-1 bg-base-100 rounded-box shadow-lg p-4">
        <h2 class="text-xl font-semibold mb-2">Mars Surface Map</h2>
        <div class="map-container w-[50vw] h-[50vw] bg-white rounded-lg overflow-hidden">
          <canvas 
            bind:this={canvas} 
            class="w-full h-full"
            aria-label="Mars surface map with 16 chunks"
          ></canvas>
        </div>
        <div class="mt-2 text-sm">
          <p>Click on a chunk to select it, then click on another chunk to swap their positions.</p>
          {#if selectedChunk !== null}
            <p class="text-green-500">Chunk selected! Click another chunk to swap.</p>
          {/if}
        </div>
      </div>
      
      <div class="flex-1 flex flex-col gap-4">
        <div class="bg-base-100 rounded-box shadow-lg p-4">
          <h2 class="text-xl font-semibold mb-2">Module Information</h2>
          <div class="overflow-x-auto">
            <table class="table table-zebra w-full">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Name</th>
                  <th>Count</th>
                </tr>
              </thead>
              <tbody>
                {#each modules as module}
                  <tr>
                    <td>{module.id}</td>
                    <td>{module.name}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>
        
        <div class="bg-base-100 rounded-box shadow-lg p-4">
          <h2 class="text-xl font-semibold mb-2">Pricing Information</h2>
          <div class="overflow-x-auto">
            <table class="table table-zebra w-full">
              <thead>
                <tr>
                  <th>Module</th>
                  <th>Cost</th>
                </tr>
              </thead>
              <tbody>
                {#each pricing as item}
                  <tr>
                    <td>{item.module}</td>
                    <td>{item.cost}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>
        
        <div class="bg-base-100 rounded-box shadow-lg p-4 mt-auto">
          <h2 class="text-xl font-semibold mb-2">Control Panel</h2>
          <div class="flex flex-wrap gap-2">
            <button class="btn btn-warning" on:click={resetGrid}>Reset Grid</button>
            <button class="btn btn-info" on:click={refreshAllChunks}>Refresh All Chunks</button>
            <div class="dropdown dropdown-top">
              <div tabindex="0" role="button" class="btn m-1">Manage URLs</div>
              <div tabindex="0" class="dropdown-content z-[1] p-2 shadow bg-base-100 rounded-box w-96">
                <div class="max-h-96 overflow-y-auto">
                  {#each chunkUrls as url, i}
                    <div class="flex items-center gap-2 mb-2">
                      <span class="text-sm font-semibold">Chunk {i}:</span>
                      <input 
                        type="text" 
                        class="input input-bordered input-sm flex-1" 
                        value={url}
                        on:change={(e) => updateChunkUrl(i, e.target.value)}
                      />
                      <button 
                        class="btn btn-xs btn-primary" 
                        on:click={() => refreshChunkData(i)}
                      >
                        Refresh
                      </button>
                    </div>
                  {/each}
                </div>
              </div>
              <!-- Add this inside the Control Panel div, after the Manage URLs dropdown -->
<div class="mt-4">
  <h3 class="text-lg font-semibold mb-2">Coordinates & Pricing</h3>
  <div class="flex items-center gap-2 mb-2">
    <input 
      type="text" 
      class="input input-bordered input-sm flex-1" 
      bind:value={coordsUrl}
      placeholder="Enter coordinates API URL"
    />
    <button 
      class="btn btn-sm btn-primary" 
      on:click={fetchCoordsData}
    >
      Update Pricing
    </button>
  </div>
  {#if listenerCoords.length > 0 && senderCoords.length > 0}
    <div class="text-sm mt-2">
      <p>Listener: [{listenerCoords.join(', ')}]</p>
      <p>Sender: [{senderCoords.join(', ')}]</p>
    </div>
  {/if}
</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</main>

<style>
  .responsive-layout {
    min-height: 60vh;
  }
  
  @media (max-width: 1023px) {
    .map-container {
      width: 100% !important;
      height: 80vw !important;
    }
  }
</style>