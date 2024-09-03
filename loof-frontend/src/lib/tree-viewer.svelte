<script lang="ts">
	import { onMount, mount } from 'svelte';
	import PocketBase from 'pocketbase';
    import type { RecordAuthResponse, RecordModel } from 'pocketbase';

	import Node from './node.svelte';
    import LoginButton from './login-button.svelte';
	import { fromDatabase, type NodeData, NodeState, toDatabase } from './node-data.svelte'

	let { pocketbase = (new PocketBase('http://127.0.0.1:8090')) }: { pocketbase: PocketBase } = $props();

	let xTargetOffset = $state(0);
	let yTargetOffset = $state(0);

	let xOffset = $state(0);
	let yOffset = $state(0);

	let targetZoom = $state(1);
	let zoom = $state(1);

	let localMousePositionX = $state(0);
	let localMousePositionY = $state(0);

	let mousePositionX = $derived((((localMousePositionX - xOffset) * targetZoom)));
	let mousePositionY = $derived((((localMousePositionY - yOffset) * targetZoom)));

	let nodes: { [id: string]: NodeData }  = $state({});

	let placingNodeId: string | null = $state(null);

	let debounceTimeout: number | null = $state(null)
	let canStopMoving: boolean = $state(false);

	function drawArrow(startX: number, startY: number, endX: number, endY: number, ctx: CanvasRenderingContext2D, targetZoom: number) {
		ctx.beginPath();

		let zoomedStartX = startX * targetZoom;
		let zoomedStartY = startY * targetZoom;

		let zoomedEndX = endX * targetZoom;
		let zoomedEndY = endY * targetZoom;

		ctx.moveTo(zoomedStartX, zoomedStartY);
		ctx.lineTo(zoomedEndX, zoomedEndY);

		let lineRads = Math.atan2(-(zoomedEndY - zoomedStartY), -(zoomedEndX - zoomedStartX));

		let leftRads = lineRads - (30 * (Math.PI / 180));
		let rightRads = lineRads + (30 * (Math.PI / 180));

		let leftX = (Math.cos(leftRads) * 25 * targetZoom) + zoomedEndX;
		let leftY = (Math.sin(leftRads) * 25 * targetZoom) + zoomedEndY;

		let rightX = (Math.cos(rightRads) * 25 * targetZoom) + zoomedEndX;
		let rightY = (Math.sin(rightRads) * 25 * targetZoom) + zoomedEndY;

		ctx.lineTo(leftX, leftY);
		ctx.moveTo(zoomedEndX, zoomedEndY);
		
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
						let newNodeData = fromDatabase(event.record, nodes[event.record.id]!.isLocal);
						nodes[event.record.id] = newNodeData;
					}
					break;
			
				default:
					break;
			}
		}, { expand: "author" });
		
		if((treeContainerDiv != null) && (treeBackground != null)) {
			let panCallbackId: number | null = null;
			let lastMouseMove: number | null = null;

			treeBackground.addEventListener("mousemove", (event) => {
				if (dragging && (event.target == event.currentTarget)) {
					const deltaTime = event.timeStamp - (lastMouseMove ?? event.timeStamp);

					if ((lastScreenX != null) && (lastScreenY != null) && (deltaTime > 0.001)) {
						xVelocity = (event.screenX - (lastScreenX)) / deltaTime;
						yVelocity = (event.screenY - (lastScreenY)) / deltaTime;

						xOffset += (event.screenX - (lastScreenX));
						yOffset += (event.screenY - (lastScreenY));

						
					}

					lastMouseMove = event.timeStamp;
					lastScreenX = event.screenX;
					lastScreenY = event.screenY;
				}
			});

			let lastPanFrame: number | null = null;

			function onNewPanFrame(timeStamp: number) {
				if (dragging == false) {
					const deltaTime = timeStamp - (lastPanFrame ?? timeStamp);

					const panDeltaX = (xVelocity) * deltaTime;
					const panDeltaY = (yVelocity) * deltaTime;

					xOffset += panDeltaX * 1;
					yOffset += panDeltaY * 1;
					
					xVelocity *= (1.0 - (0.001 * deltaTime));
					yVelocity *= (1.0 - (0.001 * deltaTime));

					if ((xVelocity != 0.00) && (yVelocity != 0.00)) {
						window.requestAnimationFrame(onNewPanFrame);

						lastPanFrame = timeStamp; 
					} else {
						lastPanFrame = null;
						panCallbackId = null;
					}
				} else {
					lastPanFrame = null;
					panCallbackId = null;
				}
			}

			let zoomCallbackId: number | null = null;

			treeBackground.addEventListener("wheel", (event) => {
				let zoomDelta = event.deltaY * 0.001;
				if ((targetZoom + zoomDelta) < 0.1) {
					zoomDelta = Math.max((targetZoom - 0.1), 0);
				} else if ((targetZoom + zoomDelta) > 3) {
					zoomDelta = Math.max((targetZoom - 3), 0);
				}
				
				targetZoom += zoomDelta;

				if (zoomCallbackId == null) {
					zoomCallbackId = window.requestAnimationFrame(onNewZoomFrame);
				}
			});

			let zoomStart: number | null = null;
			let zoomEnd: number | null = null;

			function onNewZoomFrame(timeStamp: number) {
				if (zoomStart == null) {
					zoomStart = timeStamp;
					zoomEnd = timeStamp + 50;
				} else if (zoomEnd != null) {
					zoomEnd += 50;
				}

				const deltaTime = timeStamp - zoomStart;
				const requiredZoom = targetZoom - zoom;
				const zoomDelta = (requiredZoom / (zoomEnd! - zoomStart)) * deltaTime

				xOffset -= ((localMousePositionX - xOffset) / zoom) * zoomDelta;
				yOffset -= ((localMousePositionY - yOffset) / zoom) * zoomDelta;
				zoom += zoomDelta;

				if (zoom != targetZoom) {
					window.requestAnimationFrame(onNewZoomFrame);
				} else {
					zoomStart = null;
					zoomEnd = null;

					zoomCallbackId = null;
				}
			}

			document.addEventListener("mousemove", (event) => {
				localMousePositionX = event.screenX;
				localMousePositionY = event.screenY;
				
				for (var key in nodes) {
					if ((nodes[key].state == NodeState.moving) && (nodes[key].isLocal == true)) {
						nodes[key].x = mousePositionX - 200;
						nodes[key].y = mousePositionY - 150;
					}
				}
			});

			document.addEventListener("mousedown", async (event) => {
				if ((placingNodeId != null) && (canStopMoving == true)) {
					try {
						nodes[placingNodeId].state = NodeState.editing;
						let currentNodeId = placingNodeId;

						placingNodeId = null; // if we did not copy placingNodeId to currentNodeId (and instead just used placingNodeId directly) the pocketbase await takes so long that it causes gliches (as placingNodeId should be null right afterwards (but it often isn't??!?!))

						let nodeRecord = await pocketbase.collection('nodes').update(currentNodeId, toDatabase(nodes[currentNodeId], false), { expand: "author" });	
					} catch (err) {
						if (placingNodeId != null) {
							nodes[placingNodeId].state = NodeState.moving;
						}
						console.error(err);
					}
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
				lastScreenX = null;
				lastScreenY = null;

				lastMouseMove = null;

				if (panCallbackId == null) {
					panCallbackId = window.requestAnimationFrame(onNewPanFrame);
				}
				
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
					let previousNodeXPosition = (previous_node.x! + 200);
					let previousNodeYPosition = (previous_node.y! + 200);

					drawArrow(previousNodeXPosition + (xOffset / zoom), previousNodeYPosition + (yOffset / zoom), (nodes[key].x! + 200) + (xOffset / zoom), nodes[key].y! + (yOffset / zoom), ctx, zoom);
				}
			}
		}
	})

	function endDebounce() {
		canStopMoving = true;
	}

	async function onNodeSubmission(node: NodeData) {
		if (node.isLocal == false) {
			console.warn("onNodeSubmission's node.isLocal should not be false")
		}

		if (placingNodeId == null) {
			try {
				let nodeRecord = await pocketbase.collection('nodes').create(toDatabase(node, false), { expand: "author" });
				placingNodeId = nodeRecord.id;

				nodes[nodeRecord.id] = fromDatabase(nodeRecord, true);
			} catch (err) {
				if (placingNodeId != null) {
					delete nodes[placingNodeId];
					placingNodeId = null;
				}

				console.error(err);
			}

			canStopMoving = false;
			debounceTimeout = setTimeout(endDebounce, 400);
		}
	}

	async function onNodeEndOfEditing(node: NodeData, currentText: string) {
		try {
			node.text = currentText;
			node.state = NodeState.complete;

			let record = await pocketbase.collection('nodes').update(node.id!, toDatabase(node, false));
		} catch (err) {
			node.state = NodeState.editing;

			console.error(err);
		}
	}
</script>

<div id="tree-container" style="top: {yOffset}px; left: {xOffset}px; transform: scale({zoom});">
	{#each Object.entries(nodes) as [_, node]}
		<Node bind:pocketbase={pocketbase} nodeData={node} onNodeSubmission={onNodeSubmission} onNodeEndOfEditing={onNodeEndOfEditing}></Node>
	{/each}
</div>
<canvas id="tree-background"></canvas>

<style>
	#tree-container {
		position: absolute;
	}
</style>