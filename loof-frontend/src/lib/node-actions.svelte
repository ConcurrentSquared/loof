<script lang="ts">
	import PocketBase from 'pocketbase';
    import Node from './node.svelte';
    import { mount } from 'svelte';

	import { NodeState, toDatabase, type NodeData } from './node-data.svelte'
    import Switchbox from './switchbox.svelte';

	let { bookmarks = "0", likes = "0", nodeData={id: null,

												state: NodeState.moving,

												authorId: "test",
												previousNodeId: "",

												x: 0,
												y: 0,

												text: ""}, text="", newNodeArray=$bindable([]) }: { bookmarks: string, likes: string, nodeData: NodeData, text: string, newNodeArray: Array<NodeData> } = $props();
	let currentNodeIndex: number | null = $state(null);
	let debounceTimeout: number | null = $state(null);

	let canStopMoving = $state(false);

	let isSwitchboxOpen = $state(false);
	let switchboxPositionX: number = $state(0);
	let switchboxPositionY: number = $state(0);

	let NewBranchButton: HTMLElement | null = $state(null);

	async function addNode(event: MouseEvent) {
		isSwitchboxOpen = true;
		switchboxPositionX = event.offsetX + NewBranchButton!.offsetLeft;
		switchboxPositionY = event.offsetY + NewBranchButton!.offsetTop;

		if (currentNodeIndex == null) {
			currentNodeIndex = newNodeArray.length;
			newNodeArray.push({id: null,

											state: NodeState.moving,

											authorId: "test",
											previousNodeId: nodeData.id!,

											x: 0,
											y: 0,

											text: ""});

			canStopMoving = false;
			debounceTimeout = setTimeout(endDebounce, 400);
		}
	}

	// NOTE: This function is called from the **new node**; do not use newNodeArray or currentNodeIndex to access the new node's data (because you are actually accessing the 'new new node's' data, which doesn't exist)
	async function submitNode() {
		let pocketbase = new PocketBase('http://127.0.0.1:8090');
		nodeData.authorId = "lh485oxdij1oyoa" // todo: replace
		nodeData.text = text;

		await pocketbase.collection('nodes').create(toDatabase(nodeData));

		nodeData.state = NodeState.complete;

		currentNodeIndex = null;
	}

	async function endDebounce() {
		canStopMoving = true
	}

	async function addBookmarks() {
		let pocketbase = new PocketBase('http://127.0.0.1:8090');

		let oldNodeData = await pocketbase.collection('nodes').getOne(nodeData.id!, {});
		let newNodeData = await pocketbase.collection('nodes').update(nodeData.id!, {"bookmarks": (oldNodeData.bookmarks + 1) });
		bookmarks = newNodeData.bookmarks;
	}

	async function addLikes() {
		let pocketbase = new PocketBase('http://127.0.0.1:8090');

		let oldNodeData = await pocketbase.collection('nodes').getOne(nodeData.id!, {});
		let newNodeData = await pocketbase.collection('nodes').update(nodeData.id!, {"likes": (oldNodeData.likes + 1) });
		likes = newNodeData.likes;
	}

	async function onClick(event: MouseEvent) {
		if ((currentNodeIndex != null) && (canStopMoving == true)) {
			newNodeArray[currentNodeIndex].state = NodeState.editing;

			currentNodeIndex = null;
		}
	}
</script>

<div class="node-actions-container">
	{#if nodeData.state == NodeState.moving}
	<p>Click the left mouse button to place the node</p>
	{:else if nodeData.state == NodeState.editing}
	<button onclick={submitNode}>Submit</button>
	{:else}
	<button  onclick={addNode} bind:this={NewBranchButton}>New Branch</button>
	{#if isSwitchboxOpen == true }
	<Switchbox x_position={switchboxPositionX} y_position={switchboxPositionY}></Switchbox>
	{/if}
	<button onclick={addBookmarks}>Bookmark: {bookmarks}</button>
	<button onclick={addLikes}>Like: {likes}</button>
	{/if}
</div>

<svelte:window onclick={onClick}/>

<style>
	.node-actions-container {
		position: relative;
		display: flex;

		justify-content: space-between;
    	height: 20%;

		border: 5px solid rgb(66, 66, 66);
		box-sizing: border-box;

		background-color: rgb(236, 232, 223);
	}

	button {
		background-color: inherit;

		border: 0;

  		padding: 0;
		padding-left: 5px;
		padding-right: 5px;

		height: 100%;
	}
</style>