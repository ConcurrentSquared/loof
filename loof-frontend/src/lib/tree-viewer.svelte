<script lang="ts">
	import { onMount } from 'svelte';

	let xOffset = 50;
	let yOffset = 150;

	onMount(() => {
		let dragging = false;
		let xVelocity = 0;
		let yVelocity = 0;

		let lastScreenX: number | null = null;
		let lastScreenY: number | null = null;

		let treeContainerDiv = document.getElementById("tree-container");
		if(treeContainerDiv != null) {
			treeContainerDiv.style.left = xOffset.toString() + "px";
			treeContainerDiv.style.top = yOffset.toString() + "px";

			treeContainerDiv.addEventListener("mousemove", (event) => {
				if (dragging) {
					xVelocity = event.screenX - (lastScreenX ?? xOffset);
					yVelocity = event.screenY - (lastScreenY ?? yOffset);

					console.log(xVelocity.toString())

					xOffset += xVelocity;
					yOffset += yVelocity;

					lastScreenX = event.screenX;
					lastScreenY = event.screenY;

					xVelocity = 0;
					yVelocity = 0;
				}
			});

			treeContainerDiv.addEventListener("mousedown", (event) => { 
				if (event.button == 0) {
					dragging = true;
				
					lastScreenX = event.screenX;
					lastScreenY = event.screenY;
				}
			});

			treeContainerDiv.addEventListener("mouseup", (event) => { 
				dragging = false;
			});

			treeContainerDiv.addEventListener("mouseleave", (event) => { 
				dragging = false;
			});
		}
  	});
</script>

<div id="tree-container" style="top: {yOffset}px; left: {xOffset}px">
	<slot />
</div>

<style>
	#tree-container {
		position: absolute;

		padding: 500px;
	}
</style>