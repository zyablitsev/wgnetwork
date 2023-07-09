<script>
  import { onMount } from 'svelte';
  import { get } from 'svelte/store';

  import { authToken } from '../../state';
  import Nav from './Nav.svelte';
  import LinkButton from './Button/LinkButton.svelte';
  import Spinner from './Spinner.svelte';

  export let cfg;
  let client = cfg.client;

  let session = get(authToken);
  let isDisabled = false;
  let isLoading = false;

  function handleSignOut(event) {
    event.preventDefault();

    isDisabled = true;
    isLoading = true;

    client.Fetch('auth/signout', {}, session)
      .then(result => {
        isLoading = false;
        isDisabled = false;
        authToken.set('');
      })
      .catch(err => {
        console.error(err);

        isLoading = false;
        isDisabled = false;
        authToken.set('');
      });
  }

  onMount(() => {
    let elBtnSignOut = document.getElementById('btn_signout');
    elBtnSignOut.addEventListener('click', handleSignOut);

    return () => {
      elBtnSignOut.removeEventListener('click', handleSignOut);
    }
  });
</script>

<div class="bg-white shadow">
  <div class="mx-auto max-w-4xl px-4 lg:px-8">
    <div class="flex h-16 justify-between">
      <nav class="flex space-x-6">
        <Nav {cfg} />
      </nav>

      <div class="flex items-center">
        <LinkButton id="btn_signout"
                       isDisabled={isDisabled} >
          {#if isLoading}
            <Spinner css='text-white mx-auto' />
          {:else}
            sign out
          {/if}
        </LinkButton>
      </div>
    </div>
  </div>
</div>
