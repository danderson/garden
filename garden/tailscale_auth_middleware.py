from django.contrib.auth.middleware import RemoteUserMiddleware
from django.conf import settings

class TailscaleAuthMiddleware(RemoteUserMiddleware):
    header = 'HTTP_X_TAILSCALE_USER'
    dev_user = 'danderson@github'

    def process_request(self, request):
        if self.header not in request.META and settings.DEBUG:
            request.META[self.header] = self.dev_user
        super().process_request(request)