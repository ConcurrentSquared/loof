<script lang="ts">
	let { xPosition = "0", yPosition = "0", text = "" } = $props();
</script>

<div class="node" style="top: {yPosition}px; left: {xPosition}px">
	<textarea>{text}</textarea>
	<button id="like-button">Like</button>
</div>

<style>
	.node {
		position: absolute;

		padding: 50px;

		background-color: white;

		border: 15px solid black;
		border-radius: 30px;
		box-sizing: border-box;

		width: 500px;
		height: 500px;
	}
</style>