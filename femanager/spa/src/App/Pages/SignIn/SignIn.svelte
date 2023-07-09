<script>
  import { onMount } from 'svelte';

  import { setBodyClass, setHtmlClass } from '../../util/dom';

  import { RPCError } from '../../../lib/rpcapi';
  import { authToken } from '../../state';
  import ExclamationCircle from '../../Shared/Components/Icon/Mini/ExclamationCircle.svelte';
  import PrimaryButton from '../../Shared/Components/Button/PrimaryButton.svelte';
  import Spinner from '../../Shared/Components/Spinner.svelte';

  export const params = {};
  export let cfg;
  let client = cfg.client;

  setHtmlClass('h-full', 'bg-gray-50');
  setBodyClass('h-full');

  let code = '';
  let isDisabled = true;
  let isLoading = false;

  let cssInputCodeDefault = `block w-full
 appearance-none rounded-md
 border border-gray-300
 placeholder-gray-400
 shadow-sm
 focus:border-indigo-500 focus:outline-none focus:ring-indigo-500
 disabled:cursor-not-allowed disabled:border-gray-200 disabled:bg-gray-50 disabled:text-gray-500
 text-sm touch-none`

  let cssInputCodeOnErr = `block w-full
 appearance-none rounded-md
 border border-red-300
 pr-10
 text-red-900
 placeholder-red-300
 focus:border-red-500 focus:outline-none focus:ring-red-500
 sm:text-sm touch-none`

  let cssInputCode = cssInputCodeDefault;
  let errInputCode = false;

  function handleInputCode(event) {
    errInputCode = false;
    cssInputCode = cssInputCodeDefault;

    code = event.target.value;
    if (code.length > 5 && code.length < 7) {
      isDisabled = false;
    } else {
      isDisabled = true;
    }
  }

  function handleSignIn(event) {
    event.preventDefault();

    isLoading = true;
    isDisabled = true;

    client.Fetch('auth/signin', {"code": code})
      .then(result => {
        isLoading = false;
        isDisabled = false;
        authToken.set(result.session);
      })
      .catch(err => {
        console.error(err);

        // TODO: rework
        if (err instanceof RPCError) {
          isLoading = false;
          isDisabled = true;
          errInputCode = true;
          cssInputCode = cssInputCodeOnErr;
        } else {
          isLoading = false;
          isDisabled = false;
          errInputCode = false;
        }
      });
  }

  onMount(() => {
    let elFormInputCode = document.getElementById('code');
    elFormInputCode.addEventListener('input', handleInputCode);

    let elFormBtnSubmit = document.getElementById('btn_submit');
    elFormBtnSubmit.addEventListener('click', handleSignIn);

    return () => {
      elFormBtnSubmit.removeEventListener('click', handleSignIn);
      elFormInputCode.removeEventListener('input', handleInputCode);
    }
  });
</script>

<svelte:head>
  <title>signin</title>
</svelte:head>

<div class="flex flex-col justify-center min-h-full py-12 px-8">

  <div class="mx-auto w-full max-w-md">
    <h2 class="mt-6 text-center text-3xl font-bold tracking-tight text-gray-900">
      security verification
    </h2>
  </div>

  <form class="mx-auto mt-8 w-full max-w-md space-y-6 rounded-lg bg-white py-8 px-10 shadow">

    <div>
      <label for="code"
             class="block text-sm font-medium text-gray-700">
        auth code
      </label>
      <div class="relative mt-1 rounded-md shadow-sm">
        <input id="code"
               type="password"
               name="code"
               required="true"
               class="{cssInputCode}">
        {#if errInputCode}
        <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-3">
          <ExclamationCircle />
        </div>
        {/if}
      </div>
      {#if errInputCode}
      <p class="mt-2 text-sm text-red-600">wrong code.</p>
      {/if}
    </div>

    <div class="text-center">
      <PrimaryButton type="submit"
                     id="btn_submit"
                     cssSizing='min-w-[25%]'
                     isDisabled={isDisabled} >
        {#if isLoading}
          <Spinner css='text-white mx-auto' />
        {:else}
          submit
        {/if}
      </PrimaryButton>
    </div>

  </form>

</div>
