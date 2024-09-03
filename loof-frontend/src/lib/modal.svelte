<div class="modal-container">
	<slot />
</div>

<style>
	.modal-container {
		position: fixed;

		inset: 0;

		width: fit-content;
		height: fit-content;

		margin: auto;
		padding: 30px;

		border: 10px solid black;
		background-color: white;
	}
</style>