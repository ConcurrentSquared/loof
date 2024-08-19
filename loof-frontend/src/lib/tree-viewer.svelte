<script lang="ts">
	import { onMount, mount } from 'svelte';
	import PocketBase from 'pocketbase';
    import type { RecordAuthResponse, RecordModel } from 'pocketbase';

	import Node from './node.svelte';
    import LoginButton from './login-button.svelte';
	import { fromDatabase, type NodeData, NodeState } from './node-data.svelte'

	let xOffset = $state(0);
	let yOffset = $state(0);

	let localMousePositionX = $state(0);
	let localMousePositionY = $state(0);

	let mousePositionX = $derived((localMousePositionX - xOffset));
	let mousePositionY = $derived((localMousePositionY - yOffset));

	let nodes: NodeData[] = $state([]);

	let authData: RecordAuthResponse<RecordModel> | null = $state(null);

	let inConstructionNodes: NodeData[] = $state([]);
	
	function drawArrow(startX: number, startY: number, endX: number, endY: number, ctx: CanvasRenderingContext2D) {
		ctx.beginPath();

		ctx.moveTo(startX, startY);
		ctx.lineTo(endX, endY);

		let lineRads = Math.atan2(-(endY - startY), -(endX - startX));

		let leftRads = lineRads - (30 * (Math.PI / 180));
		let rightRads = lineRads + (30 * (Math.PI / 180));

		let leftX = (Math.cos(leftRads) * 25) + endX;
		let leftY = (Math.sin(leftRads) * 25) + endY;

		let rightX = (Math.cos(rightRads) * 25) + endX;
		let rightY = (Math.sin(rightRads) * 25) + endY;

		ctx.lineTo(leftX, leftY);
		ctx.moveTo(endX, endY);
		
		ctx.lineTo(rightX, rightY);

		ctx.closePath();
		ctx.stroke();
	}

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

		let nodeRecords = await pb.collection('nodes').getFullList({});

		nodeRecords.forEach(nodeRecord => {
			nodes = [...nodes, fromDatabase(nodeRecord)];
		});
		
		if((treeContainerDiv != null) && (treeBackground != null)) {
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

			document.addEventListener("mousemove", (event) => {
				localMousePositionX = event.clientX - 200;
				localMousePositionY = event.clientY - 100;
				
				inConstructionNodes.forEach(node => {
					if (node.state == NodeState.moving) {
						node.x = mousePositionX;
						node.y = mousePositionY;
					}
				});
			})

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
		}
  	});

	$effect(() => {
		let treeContainerDiv = document.getElementById("tree-container");
		let treeBackground = document.getElementById("tree-background") as HTMLCanvasElement;
		let ctx = treeBackground.getContext("2d");

		ctx?.reset();

		if (ctx != null) {
			nodes.forEach(node => {
				let previous_node = nodes.find((databaseNode => (databaseNode.id == node.previousNodeId)));
				if (previous_node != null) {
					let previousNodeXPosition = previous_node.x! + 200;
					let previousNodeYPosition = previous_node.y! + 200;

					drawArrow(previousNodeXPosition + xOffset, previousNodeYPosition + yOffset, (node.x! + 200) + xOffset, node.y! + yOffset, ctx);
				}
			});

			inConstructionNodes.forEach(node => {
				let previous_node = nodes.find((databaseNode => (databaseNode.id == node.previousNodeId)));
				if (previous_node != null) {
					let previousNodeXPosition = previous_node.x! + 200;
					let previousNodeYPosition = previous_node.y! + 200;

					drawArrow(previousNodeXPosition + xOffset, previousNodeYPosition + yOffset, ((node.x ?? mousePositionX) + 200) + xOffset, (node.y ?? mousePositionY) + yOffset, ctx);
				}
			})
		}
	})
</script>

<div id="tree-container" style="top: {yOffset}px; left: {xOffset}px">
	{#each nodes as node}
		<Node nodeData={node} bind:newNodeArray={inConstructionNodes}></Node>
	{/each}
	{#each inConstructionNodes as constructionNode}
		<Node nodeData={constructionNode} bind:newNodeArray={inConstructionNodes}></Node>
	{/each}
</div>
<canvas id="tree-background"></canvas>
<LoginButton authData={authData}></LoginButton>

<style>
	#tree-container {
		position: absolute;
	}
</style>