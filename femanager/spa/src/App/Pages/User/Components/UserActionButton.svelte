<script>
  import { createEventDispatcher, onMount } from 'svelte';

  import PrimaryButton from '../../../Shared/Components/Button/PrimaryButton.svelte';
  import DangerButton from '../../../Shared/Components/Button/DangerButton.svelte';

  export let isDisabled = false;

  const dispatch = createEventDispatcher();

  function emitEventRemove() {
    dispatch('message', {
      action: 'remove'
    });
  }

  function emitEventEdit() {
    dispatch('message', {
      action: 'edit'
    });
  }

  let btnRemoveId = 'btn_remove';
  let btnEditId = 'btn_edit'

  onMount(() => {
    let elBtnRemove = document.getElementById(btnRemoveId);
    elBtnRemove.addEventListener('click', emitEventRemove);

    let elBtnEdit = document.getElementById(btnEditId);
    elBtnEdit.addEventListener('click', emitEventEdit);

    return () => {
      elBtnRemove.removeEventListener('click', emitEventRemove);
      elBtnEdit.removeEventListener('click', emitEventEdit);
    }
  });
</script>

<div class="mx-auto my-5 max-w-3xl bg-white shadow sm:rounded-lg">

  <div class="border-b border-gray-200 px-4 py-5 sm:px-6">
    <div class="flex justify-between">

      <DangerButton
        {isDisabled}
        id={btnRemoveId}
        cssSizing='w-32'
        cssFlex='inline-flex justify-center'
        cssSpacing='py-2 px-4'>remove</DangerButton>

      <PrimaryButton
        {isDisabled}
        id={btnEditId}
        cssSizing='w-32'
        cssFlex='inline-flex justify-center'
        cssSpacing='py-2 px-4'>edit</PrimaryButton>

    </div>
  </div>

</div>
