<script>
  import { onMount } from 'svelte';
  import QRCodeStyling from 'qr-code-styling';

  export let text;
  export let id;

  onMount(() => {
    const qrCode = new QRCodeStyling({
        width: 250,
        height: 250,
        type: "svg",
        data: text,
        dotsOptions: {
            color: "#000000",
            type: "rounded"
        },
        backgroundOptions: {
            color: "#ffffff",
        },
        imageOptions: {
            crossOrigin: "anonymous",
            margin: 20
        }
    });

    qrCode.append(document.getElementById(id));
  });
</script>

<slot></slot>
