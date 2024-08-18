<script lang="ts">
	import PocketBase from 'pocketbase';
    import Node from './node.svelte';
    import { mount } from 'svelte';

	import type { InConstructionNode } from './node.ts'

	let { bookmarks = "0", likes = "0", nodeId="", newNodeArray=[] } = $props();

	async function addNode() {
		let arr = newNodeArray as Array<InConstructionNode>;
		arr.push({
			x: 0,
			y: 0
		})
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
</script>

<div class="node-actions-container">
	<button onclick={addNode}>New Branch</button>

	<button onclick={addBookmarks}>Bookmark: {bookmarks}</button>
	<button onclick={addLikes}>Like: {likes}</button>
</div>

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