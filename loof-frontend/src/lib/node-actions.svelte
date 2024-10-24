<script lang="ts">
	import PocketBase from 'pocketbase';

	import { checkIfUserIsAuthor, checkIfUserIsAuthorOfNode, fromDatabase, getDisplayUsername, getUserAuthor, NodeState, toDatabase, type NodeData } from './node-data.svelte'
    import Switchbox from './switchbox.svelte';
    import { mount, onMount, tick } from 'svelte';

	let { pocketbase = $bindable(new PocketBase('http://127.0.0.1:8090')), isLogin = ((pocketbase.authStore.model == null) ? false : true), bookmarks = "0", nodeData= {id: null,
												
												state: NodeState.moving,
												isLocal: true,
												fromRobot: false,

												authorId: "test",
												previousNodeId: "",

												x: 0,
												y: 0,

												text: ""}, text="", onNodeSubmission, onNodeEndOfEditing, updatePocketbaseRequest }: { pocketbase: PocketBase, isLogin: boolean, bookmarks: string, nodeData: NodeData, text: string, onNodeSubmission: (node: NodeData) => void, onNodeEndOfEditing: (node: NodeData, currentText: string) => void, updatePocketbaseRequest: (node: NodeData, currentText: string) => void } = $props();

	let isSwitchboxOpen = $state(false);
	let switchboxPositionX: number = $state(0);
	let switchboxPositionY: number = $state(0);

	let NewBranchButton: HTMLElement | null = $state(null);

	let authorUsername: string | null = $state(null);
	let fromCurrentUser: boolean = $state(false)

	let updateInterval: number | null = $state(null);

	async function openSwitchbox(event: MouseEvent) {
		isSwitchboxOpen = true;
		switchboxPositionX = event.offsetX + NewBranchButton!.offsetLeft;
		switchboxPositionY = event.offsetY + NewBranchButton!.offsetTop;
	}
	
	async function addHumanNode() {
		isSwitchboxOpen = false;

		if (pocketbase.authStore.model != null) {
			let authorRecord = await pocketbase.collection('authors').getFirstListItem('author_id?="' + pocketbase.authStore.model.id + '"&&' + 'origin="human"')
							   .catch(err => {return pocketbase.collection('authors').create({ "author_id": pocketbase.authStore.model!.id, "origin": "human"})})
			
			let newNodeData: NodeData = {id: null,
							state: NodeState.moving,
							isLocal: true,
							fromRobot: false,

							authorId: authorRecord.id,
							previousNodeId: nodeData.id!,
		
							x: 0,
							y: 0,

							text: ""}
			
			onNodeSubmission(newNodeData);
		}
	}

	async function addAINode() {
		isSwitchboxOpen = false;

		if (pocketbase.authStore.model != null) {
			let newNodeData: NodeData = {id: null,
							state: NodeState.moving,
							isLocal: true,
							fromRobot: false,

							authorId: "lh485oxdij1oyoa",
							previousNodeId: nodeData.id!,
		
							x: 0,
							y: 0,

							text: ""}

			onNodeSubmission(newNodeData);
		}
	}

	// NOTE: This function is called from the **new node**; do not use newNodeArray or currentNodeIndex to access the new node's data (because you are actually accessing the 'new new node's' data, which doesn't exist)
	async function submitNode() {
		onNodeEndOfEditing(nodeData, text)
	}

	async function addBookmarks() {
		try {
			let res = await pocketbase.collection('bookmarks').create({ "node": nodeData.id!, "user": pocketbase.authStore.model!.id });
			bookmarks += 1;
		} catch (error) {
			console.error(error);
		}
	}

	async function addReport() {
		try {
			let res = await pocketbase.collection('bookmarks').create({ "node": nodeData.id!, "user": pocketbase.authStore.model!.id });
			//bookmarks += 1;
		} catch (error) {
			console.error(error);
		}
	}

	async function onMouseDown(event: MouseEvent) {
		if ((isSwitchboxOpen == true) && (event.target == document.getElementById("tree-background"))) {
			isSwitchboxOpen = false;
		}
	}

	export async function updateAuthorship(userId: string)
	{
		let currentAuthorId = await getUserAuthor(userId, pocketbase);
		fromCurrentUser = await checkIfUserIsAuthorOfNode(currentAuthorId, pocketbase);
	}

	onMount(async () => {
		authorUsername = await getDisplayUsername(nodeData, pocketbase);

		let currentUserId: string | null = null;
		let currentUserStore = pocketbase.authStore;
		if (currentUserStore.model != null) {
			currentUserId = currentUserStore.model.id;
		}

		let currentAuthorId = await getUserAuthor(currentUserId, pocketbase);
		fromCurrentUser = await checkIfUserIsAuthorOfNode(currentAuthorId, pocketbase);

		bookmarks = (await pocketbase.collection('most_bookmarked').getOne(nodeData.id!)).bookmark_count;

		pocketbase.collection('most_bookmarked').subscribe(nodeData.id!, (event) => {
			if (event.action == "update") {
				bookmarks = event.record.bookmark_count;
			}
		});

		updateInterval = setInterval(async () => { updatePocketbaseRequest(nodeData, text); }, 5000);
	});
</script>

<div class="node-actions-container">
	{#if nodeData.state == NodeState.moving}
	<p>Click the left mouse button (or tap) to place the node</p>
	{:else if (nodeData.state == NodeState.editing) && (nodeData.fromRobot == false) && (fromCurrentUser == true)}
	<button onclick={submitNode}>Submit</button>
	{:else if (nodeData.state == NodeState.editing) && (nodeData.fromRobot == false) && (fromCurrentUser == false)}
	<p>User {authorUsername} is writing..</p>
	{:else if (nodeData.state == NodeState.editing) && (nodeData.fromRobot == true)}
	<p>Generating...</p>
	{:else}
	<button onclick={openSwitchbox} bind:this={NewBranchButton} disabled={!isLogin}>New Branch</button>
	{#if isSwitchboxOpen == true }
	<Switchbox x_position={switchboxPositionX} y_position={switchboxPositionY} on_write_selection={addHumanNode} on_generate_selection={addAINode}></Switchbox>
	{/if}
	<button onclick={addBookmarks} disabled={!isLogin}>Bookmark: {bookmarks}</button>
	<button onclick={addReport} disabled={!isLogin}>Report</button>
	<p>Written by {authorUsername}</p>
	{/if}
</div>

<svelte:window on:mousedown={onMouseDown}/>

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
		font: 0.85em "Jost", Arial, Helvetica, sans-serif;
		
		background-color: inherit;

		border: 0;

  		padding: 0;
		padding-left: 5px;
		padding-right: 5px;

		height: 100%;

		cursor: pointer;
	}

	button:disabled {
		color: #646464;
		cursor: not-allowed;
	}

	p {
		font: 0.85em "Jost", Arial, Helvetica, sans-serif;

		padding: 0;
		padding-left: 5px;
		padding-right: 5px;

		margin: 0;
		margin-top: auto;
		margin-bottom: auto;
	}
</style>