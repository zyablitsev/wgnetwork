<script>
  import { onMount } from 'svelte';
  import { get } from 'svelte/store';

  import { setBodyClass, setHtmlClass } from '../../util/dom';
  import { ValidationError, RPCError } from '../../../lib/rpcapi';
  import { getUser, formValidate } from './func.js';
  import { authToken, moveBack, moveTo } from '../../state';

  import PrimaryButton from '../../Shared/Components/Button/PrimaryButton.svelte';
  import LinkButton from '../../Shared/Components/Button/LinkButton.svelte';
  import CardHeading from '../../Shared/Components/CardHeading.svelte';
  import Header from '../../Shared/Components/Header.svelte';
  import Spinner from '../../Shared/Components/Spinner.svelte';
  import Toggle from '../../Shared/Components/Toggle.svelte';

  export let params = {};
  export let cfg;
  let client = cfg.client;
  let router = cfg.router;

  setHtmlClass('h-full', 'bg-gray-100');
  setBodyClass('h-full');

  let session = get(authToken);
  let isDisabled = true;
  let isLoading = true;

  let user = {};
  let title = 'edit user';

  let isManager = false;
  let username = '';

  function formHandleUsername(event) {
    event.preventDefault;

    username = event.target.value;
    isDisabled = !formValidate(username);
  }

  function formToggleIsManager(event) {
    event.preventDefault;

    isManager = event.detail.isChecked;
    isDisabled = !formValidate(username);
  }

  function handleUserEdit(event) {
    event.preventDefault();

    isLoading = true;
    isDisabled = true;

    let params = {'uuid': user.uuid, 'name': username, 'is_manager': isManager};
    client.Fetch('manager/user/edit', params, session)
      .then(result => {
        isLoading = false;
        isDisabled = false;
        let path = router.ReverseURI('user', {'uuid': result.uuid});
        moveTo(path);
      })
      .catch(err => {
        console.error(err);

        // TODO: rework
        if (err instanceof RPCError) {
          isLoading = false;
          isDisabled = true;
        } else {
          isLoading = false;
          isDisabled = false;
          console.error(err);
        }
      });
  }

  onMount(() => {
    getUser(client, session, params['uuid'])
      .then(result => {
        authToken.set(result['session']);

        user = result['user'];
        username = user['name'];
        isManager = user['isManager'];
        title = 'edit user - ' + user['name'];
        isLoading = false;
        isDisabled = !formValidate(username);
      })
      .catch(err => {
        console.error(err);

        // TODO: rework
        if (err instanceof RPCError) {
          isLoading = false;
        } else {
          isLoading = false;
        }
      });

    let elInputUsername = document.getElementById('username');
    elInputUsername.addEventListener('input', formHandleUsername);
    let elBtnCancel = document.getElementById('btn_cancel');
    elBtnCancel.addEventListener('click', moveBack);
    let elBtnSave = document.getElementById('btn_submit');
    elBtnSave.addEventListener('click', handleUserEdit);

    return () => {
      elBtnSave.removeEventListener('click', handleUserEdit);
      elBtnCancel.removeEventListener('click', moveBack);
      elInputUsername.removeEventListener('input', formHandleUsername);
    }
  });
</script>

<svelte:head>
  <title>{title}</title>
</svelte:head>

<Header {cfg} />

<div class="mx-auto my-5 max-w-3xl bg-white shadow sm:rounded-lg">

  <CardHeading {isLoading} title='edit user' description='the user with the manager role is allowed access to the control panel' />

  <div class="p-0">
    <form>

      <dl class="divide-y divide-gray-200">

        <div class="grid grid-cols-5 gap-4 py-5 px-6">
          <dt class="text-sm font-medium text-gray-500">username:</dt>
          <dd class="col-span-4 mt-0 text-sm text-gray-900">
            <input type="text"
                   name="username"
                   id="username"
                   class="block w-full max-w-lg rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:max-w-xs sm:text-sm touch-none" value={user.name}>
          </dd>
        </div>

        <div class="grid grid-cols-5 gap-4 py-5 px-6">
          <dt class="text-sm font-medium text-gray-500">is manager:</dt>
          <dd class="col-span-4 mt-0 text-sm text-gray-900">
            {#if isManager}
              <Toggle isChecked={true}
                      name='is_manager' id='is_manager'
                      on:message={formToggleIsManager} />
            {:else}
              <Toggle isChecked={false}
                      name='is_manager' id='is_manager'
                      on:message={formToggleIsManager} />
            {/if}
          </dd>
        </div>

      </dl>

      <div class="flex justify-end py-5 px-6">

        <LinkButton type='button'
                    id='btn_cancel'
                    cssFlex=''
                    cssSpacing='py-2 px-4'>
            cancel
        </LinkButton>

        <PrimaryButton type='submit'
                       id='btn_submit'
                       cssFlex='inline-flex justify-center'
                       cssSpacing='py-2 px-4 ml-3'
                       isDisabled={isDisabled}>
          {#if isLoading}
            <div class="inline-flex items-center object-center">
              <Spinner css='text-white' />
              processing
            </div>
          {:else}
            save
          {/if}
        </PrimaryButton>

      </div>

    </form>
  </div>

</div>
