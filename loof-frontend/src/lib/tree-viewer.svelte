<script lang="ts">
	import { onMount, mount } from 'svelte';
	import PocketBase from 'pocketbase';
    import type { RecordModel } from 'pocketbase';

	import Node from './node.svelte';

	let xOffset = 50;
	let yOffset = 150;

	let node: RecordModel | null = null;

	onMount(async () => {
		let dragging = false;
		let xVelocity = 0;
		let yVelocity = 0;

		let lastScreenX: number | null = null;
		let lastScreenY: number | null = null;

		let treeContainerDiv = document.getElementById("tree-container");
		let treeBackground = document.getElementById("tree-background") as HTMLCanvasElement;
		let ctx = treeBackground.getContext("2d");

		let pb = new PocketBase('http://127.0.0.1:8090');

		let nodeData = await pb.collection('nodes').getFullList({});

		node = nodeData[0]
		
		if((treeContainerDiv != null) && (treeBackground != null) && (ctx != null)) {
			treeContainerDiv.style.left = xOffset.toString() + "px";
			treeContainerDiv.style.top = yOffset.toString() + "px";

			treeBackground.addEventListener("mousemove", (event) => {
				if (dragging && (event.target == event.currentTarget)) {
					xVelocity = event.screenX - (lastScreenX ?? xOffset);
					yVelocity = event.screenY - (lastScreenY ?? yOffset);

					xOffset += xVelocity;
					yOffset += yVelocity;

					lastScreenX = event.screenX;
					lastScreenY = event.screenY;

					xVelocity = 0;
					yVelocity = 0;
				}
			});

			treeBackground.addEventListener("mousedown", (event) => { 
				if (event.button == 0) {
					dragging = true;
				
					lastScreenX = event.screenX;
					lastScreenY = event.screenY;
				}
			});

			treeBackground.addEventListener("mouseup", (event) => { 
				dragging = false;
			});

			treeBackground.addEventListener("mouseleave", (event) => { 
				dragging = false;
			});

			treeBackground.width = window.innerWidth;
			treeBackground.height = window.innerHeight;
			
			ctx.beginPath();

			ctx.moveTo(0, 0);
			ctx.lineTo(50, 50);

			let lineRads = Math.atan2(-50, -50);

			let leftRads = lineRads - (30 * (Math.PI / 180));
			let rightRads = lineRads + (30 * (Math.PI / 180));

			let leftX = (Math.cos(leftRads) * 25) + 50;
			let leftY = (Math.sin(leftRads) * 25) + 50;

			let rightX = (Math.cos(rightRads) * 25) + 50;
			let rightY = (Math.sin(rightRads) * 25) + 50;

			ctx.lineTo(leftX, leftY);
			ctx.moveTo(50, 50);
			
			ctx.lineTo(rightX, rightY);

			ctx.closePath();
			ctx.stroke();
		}
  	});
</script>

<div id="tree-container" style="top: {yOffset}px; left: {xOffset}px">
	<Node xPosition={node?.x_position} yPosition={node?.y_position} text={node?.text}></Node>
</div>
<canvas id="tree-background"></canvas>

<style>
	#tree-container {
		position: absolute;
	}
</style>