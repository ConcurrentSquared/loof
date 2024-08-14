<script lang="ts">
	export let xPosition = "0";
	export let yPosition = "0";

	export let text = "";
</script>

<div class="node" style="top: {xPosition}px; left: {yPosition}px">
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
	}
</style>