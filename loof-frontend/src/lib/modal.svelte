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
		padding: 5px;

		border: 5px solid black;
	}
</style>