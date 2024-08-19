<script lang="ts">
	import PocketBase from 'pocketbase';
    import Node from './node.svelte';
    import { mount } from 'svelte';

	import { NodeState, type InConstructionNode } from './node-data.svelte'

	let { bookmarks = "0", likes = "0", nodeId="", newNodeArray=$bindable([]) } = $props();
	let currentNodeIndex: number | null = $state(null);
	let debounceTimeout: number | null = $state(null);

	let canStopMoving = $state(false);

	async function addNode() {
		if (currentNodeIndex == null) {
			let arr = newNodeArray as Array<InConstructionNode>;
			currentNodeIndex = arr.push({
				state: NodeState.moving, 
				previousNodeId: nodeId,

				x: null,
				y: null
			}) - 1;

			canStopMoving = false;
			debounceTimeout = setTimeout(endDebounce, 400);
		}
	}

	async function endDebounce() {
		canStopMoving = true
	}

	async function addBookmarks() {
		let pocketbase = new PocketBase('http://127.0.0.1:8090');

		let oldNodeData = await pocketbase.collection('nodes').getOne(nodeId, {});
		let newNodeData = await pocketbase.collection('nodes').update(nodeId, {"bookmarks": (oldNodeData.bookmarks + 1) });
		bookmarks = newNodeData.bookmarks;
	}

	async function addLikes() {
		let pocketbase = new PocketBase('http://127.0.0.1:8090');

		let oldNodeData = await pocketbase.collection('nodes').getOne(nodeId, {});
		let newNodeData = await pocketbase.collection('nodes').update(nodeId, {"likes": (oldNodeData.likes + 1) });
		likes = newNodeData.likes;
	}

	async function onClick(event: MouseEvent) {
		if ((currentNodeIndex != null) && (canStopMoving == true)) {
			let arr = newNodeArray as Array<InConstructionNode>;
			arr[currentNodeIndex].state = NodeState.editing;

			currentNodeIndex = null;
		}
	}
</script>

<div class="node-actions-container">
	<button id="new_branch_button" onclick={addNode}>New Branch</button>

	<button onclick={addBookmarks}>Bookmark: {bookmarks}</button>
	<button onclick={addLikes}>Like: {likes}</button>
</div>
<svelte:window onclick={onClick}/>

<style>
	.node-actions-container {
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