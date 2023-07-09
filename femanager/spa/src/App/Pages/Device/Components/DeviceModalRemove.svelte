<script>
  import { createEventDispatcher, onMount } from 'svelte';

  import { transitionCfg } from '../../../util/transition.js';

  import ExclamationTriangle from '../../../Shared/Components/Icon/Outline/ExclamationTriangle.svelte';
  import DangerButton from '../../../Shared/Components/Button/DangerButton.svelte';
  import LinkButton from '../../../Shared/Components/Button/LinkButton.svelte';

  export let isDisabled = false;

  const dispatch = createEventDispatcher();

  function emitEventRemove() {
    dispatch('message', {
      action: 'remove'
    });
  }

  function emitEventCancel() {
    dispatch('message', {
      action: 'cancel'
    });
  }

  let btnRemoveId = 'btn_removemodal';
  let btnCancelId = 'btn_cancelmodal'

  onMount(() => {
    let elBtnRemove = document.getElementById(btnRemoveId);
    elBtnRemove.addEventListener('click', emitEventRemove);

    let elBtnClose = document.getElementById(btnCancelId);
    elBtnClose.addEventListener('click', emitEventCancel);

    return () => {
      elBtnRemove.removeEventListener('click', emitEventRemove);
      elBtnClose.removeEventListener('click', emitEventCancel);
    }
  });
</script>

<div class="relative z-10" aria-labelledby="modal-title" role="dialog" aria-modal="true">

  <div
    in:transitionCfg|local="{{
        duration: 300,
        base: 'transition ease-out duration-300',
        from: 'opacity-0',
        to: 'opacity-100',
        isOut: false,
      }}"
    out:transitionCfg|local="{{
        duration: 200,
        base: 'transition ease-in duration-200',
        from: 'opacity-100',
        to: 'opacity-0',
        isOut: true,
      }}">
  <!--
    Background backdrop, show/hide based on modal state.

    Entering: "ease-out duration-300"
      From: "opacity-0"
      To: "opacity-100"
    Leaving: "ease-in duration-200"
      From: "opacity-100"
      To: "opacity-0"
  -->
    <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"></div>

    <div class="fixed inset-0 z-10 overflow-y-auto">
      <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">

        <div
          in:transitionCfg|local="{{
              duration: 300,
              base: 'transition ease-out duration-300',
              from: 'opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95',
              to: 'opacity-100 translate-y-0 sm:scale-100',
              isOut: false,
            }}"
          out:transitionCfg|local="{{
              duration: 200,
              base: 'transition ease-in duration-200',
              from: 'opacity-100 translate-y-0 sm:scale-100',
              to: 'opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95',
              isOut: true,
            }}">
              <!--
                Modal panel, show/hide based on modal state.

                Entering: "ease-out duration-300"
                  From: "opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
                  To: "opacity-100 translate-y-0 sm:scale-100"
                Leaving: "ease-in duration-200"
                  From: "opacity-100 translate-y-0 sm:scale-100"
                  To: "opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
              -->


        </div>

        <div class="relative transform overflow-hidden rounded-lg bg-white px-4 pt-5 pb-4 text-left shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg sm:p-6">
          <div class="sm:flex sm:items-start">

            <div class="mx-auto flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-full bg-red-100 sm:mx-0 sm:h-10 sm:w-10">
              <ExclamationTriangle />
            </div>

            <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
              <h3 class="text-lg font-medium leading-6 text-gray-900" id="modal-title">Remove device</h3>
              <div class="mt-2">
                <p class="text-sm text-gray-500">
  Are you sure you want to remove this device? All data will be permanently removed forever. This action cannot be undone.
                </p>
              </div>
            </div>

          </div>
          <div class="mt-5 sm:mt-4 sm:flex sm:flex-row-reverse">

            <DangerButton
              {isDisabled}
              id={btnRemoveId}
              cssFlex='inline-flex justify-center'
              cssSizing='w-full sm:w-auto'
              cssSpacing='py-2 px-4 sm:ml-3'
              cssTypography='text-base sm:text-sm font-medium'>remove</DangerButton>

            <LinkButton
              {isDisabled}
              id={btnCancelId}
              cssFlex='mt-3 inline-flex justify-center'
              cssSizing='w-full sm:w-auto'
              cssSpacing='py-2 px-4 sm:mt-0'
              cssTypography='text-base sm:text-sm font-medium'>cancel</LinkButton>

          </div>
        </div>
      </div>
    </div>

  </div>

</div>
