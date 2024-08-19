<script lang="ts">
    import NodeActions from "./node-actions.svelte";
    import { NodeState, type NodeData } from "./node-data.svelte";

	let { nodeData = {id: null,

					state: NodeState.moving,

					authorId: "test",
					previousNodeId: "",

					x: 0,
					y: 0,

					text: ""}, newNodeArray=$bindable([]) }: { nodeData: NodeData, newNodeArray: Array<NodeData> } = $props();

	function checkEditable(nodeData: NodeData): boolean {
		switch (nodeData.state) {
			case NodeState.editing:
				return true;
			default:
				return false;
		}
	}
	let isEditable = $derived(checkEditable(nodeData));
	let currentText = $state(nodeData.text);
</script>

<div class="node" style="top: {nodeData.y!.toString()}px; left: {nodeData.x!.toString()}px">
	<textarea class="node-text-area" readonly={!isEditable} bind:value={currentText}></textarea>
	<NodeActions bookmarks=0 likes=0 nodeData={nodeData} text={currentText} bind:newNodeArray={newNodeArray}></NodeActions>
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

		overflow: hidden;

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