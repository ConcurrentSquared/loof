<script lang="ts">
	import PocketBase from "pocketbase"

    import NodeActions from "./node-actions.svelte";
    import { NodeState, type NodeData } from "./node-data.svelte";

	let { pocketbase = $bindable(new PocketBase('http://127.0.0.1:8090')), isLogin = ((pocketbase.authStore.model == null) ? false : true), nodeData = {id: null,

					state: NodeState.moving,
					isLocal: true,
					fromRobot: false,

					authorId: "test",
					previousNodeId: "",

					x: 0,
					y: 0,

					text: ""}, onNodeSubmission, onNodeEndOfEditing, updatePocketbaseRequest }: { pocketbase: PocketBase, isLogin: boolean, nodeData: NodeData, onNodeSubmission: (node: NodeData) => void, onNodeEndOfEditing: (node: NodeData, currentText: string) => void, updatePocketbaseRequest: (node: NodeData, currentText: string) => void } = $props();

	function checkEditable(nodeData: NodeData): boolean {
		switch (nodeData.state) {
			case NodeState.editing:
				if ((nodeData.fromRobot == false) && (isLogin == true)) {
					return true;
				} else {
					return false;
				}
			default:
				return false;
		}
	}
	let isEditable = $derived(checkEditable(nodeData));
	let currentText = $derived(nodeData.text);

	let editableText = $state(nodeData.text)

	async function internalUpdate(nodeData: NodeData, text: string) {
		updatePocketbaseRequest(nodeData, text);
		editableText = nodeData.text;
	}

	let display = $state((nodeData.y! > -1000) && (nodeData.x! > -1000) && (nodeData.y! < 1000) && (nodeData.x! < 1000) ? 'block' : 'none')

	let nodeActions: any | null = $state(null);
	export async function updateAuthorship(userId: string) {
		console.log(userId)
		await nodeActions?.updateAuthorship(userId);
	}
</script>

<div class="node" style="top: {nodeData.y!.toString()}px; left: {nodeData.x!.toString()}px; display: {display};">
	{#if isEditable}
	<textarea class="node-text-area" readonly={!isEditable} bind:value={editableText}></textarea>
	{:else}
	<textarea class="node-text-area not-editable" readonly={!isEditable} value={currentText}></textarea>
	{/if}
	<NodeActions bind:this={nodeActions} bind:pocketbase={pocketbase} isLogin={isLogin} bookmarks=0 nodeData={nodeData} text={editableText} onNodeSubmission={onNodeSubmission} onNodeEndOfEditing={onNodeEndOfEditing} updatePocketbaseRequest={internalUpdate}></NodeActions>
</div>

<style>
	.node {
		display: flex;
		flex-direction: column;
		justify-content: space-between;

		position: absolute;

		background-color: white;

		box-sizing: content-box;

		width: 400px;
		height: 200px;

		overflow: visible;

		border-radius: 2.5px;
	}

	.node-text-area {
		height: 80%;
    	width: 100%;

		resize: none;

		border: 5px solid rgb(133, 133, 133);
		border-bottom: none;
		box-sizing: border-box;

		background-color: rgb(255, 252, 245);

		font-family: "Courier Prime", "Courier New", Courier, monospace;
	}

	.node-text-area.not-editable {
		cursor: default;
	}

	.node-text-area:focus {
		outline: none;

		border: 5px solid rgb(133, 133, 133);
	}

	@-moz-document url-prefix()  {
		.node-text-area {
			scrollbar-width: auto;
			scrollbar-color: rgb(100, 100, 100);
		}
	}

	.node-text-area::-webkit-scrollbar {
		display: block;

		width: 1rem;
	}

	.node-text-area::-webkit-scrollbar-track {
		background: transparent;
	}

	.node-text-area::-webkit-scrollbar-thumb {
		background: rgb(100, 100, 100);
		border-radius: 100px;
		
		border: 5px solid rgba(0, 0, 0, 0);
		background-clip: padding-box;
	}
</style>