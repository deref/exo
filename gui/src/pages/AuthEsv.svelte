<script lang="ts">
  import { api } from '../lib/api';

  /*const esvHost = 'http://localhost:5000';*/
  const esvHost = 'https://secrets.deref.io';

  (async () => {
    const queryParams = new URLSearchParams(window.location.search);
    const returnTo = queryParams.get('returnTo') ?? '/';
    const refreshToken = queryParams.get('refreshToken');
    if (refreshToken) {
      await api.kernel.saveEsvRefreshToken(esvHost, refreshToken);
      const uri = new URL(window.location.href);
      uri.search = '';
      uri.hash = returnTo;
      window.location.replace(uri.href);
      return;
    }

    const params = new URLSearchParams({
      port: window.location.port,
      exoLocation: returnTo,
    });
    const dest = `${esvHost}/api/auth/login?returnTo=${encodeURIComponent(
      `/api/send-to-exo?${params.toString()}`,
    )}`;
    window.location.href = dest;
  })();
</script>
