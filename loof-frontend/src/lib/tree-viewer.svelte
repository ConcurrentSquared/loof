<script lang="ts">
	import { onMount, mount } from 'svelte';
	import PocketBase from 'pocketbase';
    import type { RecordAuthResponse, RecordModel } from 'pocketbase';

	import Node from './node.svelte';
    import LoginButton from './login-button.svelte';
	import { fromDatabase, type NodeData, NodeState, toDatabase } from './node-data.svelte'

	let xOffset = $state(0);
	let yOffset = $state(0);

	let localMousePositionX = $state(0);
	let localMousePositionY = $state(0);

	let mousePositionX = $derived((localMousePositionX - xOffset));
	let mousePositionY = $derived((localMousePositionY - yOffset));

	let nodes: { [id: string]: NodeData }  = $state({});

	let pocketbase = $state(new PocketBase('http://127.0.0.1:8090'))

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

		let nodeRecords = await pocketbase.collection('nodes').getFullList({ expand: "author" });

		nodeRecords.forEach(nodeRecord => {
			let node = fromDatabase(nodeRecord, false);
			nodes[node.id!] = node; 
		});

		pocketbase.collection('nodes').subscribe('*', function (event) {
			switch (event.action) {
				case "create":
					if (nodes[event.record.id] == null) { 
						nodes[event.record.id] = fromDatabase(event.record, false);
					}
					break;
				case "update":
					if ((nodes[event.record.id]!.isLocal == false) || ((nodes[event.record.id]!.fromRobot == true) && (nodes[event.record.id]!.state != NodeState.moving))) {
						nodes[event.record.id] = fromDatabase(event.record, nodes[event.record.id]!.isLocal);
					}
					break;
			
				default:
					break;
			}
		}, { expand: "author" });
		
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
				
				for (var key in nodes) {
					if ((nodes[key].state == NodeState.moving) && (nodes[key].isLocal == true)) {
						nodes[key].x = mousePositionX;
						nodes[key].y = mousePositionY;
					}
				}
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
			for (var key in nodes) {
				let previous_node = nodes[nodes[key].previousNodeId]!;
				if (previous_node != null) {
					let previousNodeXPosition = previous_node.x! + 200;
					let previousNodeYPosition = previous_node.y! + 200;

					drawArrow(previousNodeXPosition + xOffset, previousNodeYPosition + yOffset, (nodes[key].x! + 200) + xOffset, nodes[key].y! + yOffset, ctx);
				}
			}
		}
	})

	async function onNodeSubmission(node: NodeData) {
		if (node.isLocal == false) {
			console.warn("onNodeSubmission's node.isLocal should not be false")
		}

		let nodeRecord = await pocketbase.collection('nodes').create(toDatabase(node, false));
		nodes[nodeRecord.id] = fromDatabase(nodeRecord, true);
		//if (currentNodeIndex == null) {
		//	currentNodeIndex = newNodeArray.length;
		//	newNodeArray.push(fromDatabase(nodeRecord, true));
		//
		//	canStopMoving = false;
		//	debounceTimeout = setTimeout(endDebounce, 400);
		//}
	}
</script>

<div id="tree-container" style="top: {yOffset}px; left: {xOffset}px">
	{#each Object.entries(nodes) as [_, node]}
		<Node bind:pocketbase={pocketbase} nodeData={node} onNodeSubmission={onNodeSubmission}></Node>
	{/each}
</div>
<canvas id="tree-background"></canvas>
<LoginButton bind:pocketbase={pocketbase}></LoginButton>

<style>
	#tree-container {
		position: absolute;
	}
</style>