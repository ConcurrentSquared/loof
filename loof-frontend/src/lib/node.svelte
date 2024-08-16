<script lang="ts">
    import NodeActions from "./node-actions.svelte";

	let { xPosition = "0", yPosition = "0", text = "", bookmarks = "0", likes = "0", nodeId = "" } = $props();
</script>

<div class="node" style="top: {yPosition}px; left: {xPosition}px">
	<textarea class="node-text-area">{text}</textarea>
	<NodeActions bookmarks={bookmarks} likes={likes} nodeId={nodeId}></NodeActions>
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