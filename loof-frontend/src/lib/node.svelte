<script lang="ts">
	import PocketBase from "pocketbase"

    import NodeActions from "./node-actions.svelte";
    import { NodeState, type NodeData } from "./node-data.svelte";

	let { pocketbase = $bindable(new PocketBase('http://127.0.0.1:8090')), nodeData = {id: null,

					state: NodeState.moving,
					isLocal: true,
					fromRobot: false,

					authorId: "test",
					previousNodeId: "",

					x: 0,
					y: 0,

					text: ""}, onNodeSubmission, onNodeEndOfEditing }: { pocketbase: PocketBase, nodeData: NodeData, onNodeSubmission: (node: NodeData) => void, onNodeEndOfEditing: (node: NodeData, currentText: string) => void } = $props();

	function checkEditable(nodeData: NodeData): boolean {
		switch (nodeData.state) {
			case NodeState.editing:
				if (nodeData.fromRobot == false) {
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
</script>

<div class="node" style="top: {nodeData.y!.toString()}px; left: {nodeData.x!.toString()}px">
	{#if isEditable}
	<textarea class="node-text-area" readonly={!isEditable} bind:value={editableText}></textarea>
	{:else}
	<textarea class="node-text-area" readonly={!isEditable} value={currentText}></textarea>
	{/if}
	<NodeActions bind:pocketbase={pocketbase} bookmarks=0 likes=0 nodeData={nodeData} text={editableText} onNodeSubmission={onNodeSubmission} onNodeEndOfEditing={onNodeEndOfEditing}></NodeActions>
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
	}

	.node-text-area:focus {
		outline: none;

		border: 5px solid rgb(133, 133, 133);
	}
</style>