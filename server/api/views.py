from django.shortcuts import render
from django.contrib.auth.models import User
from django.http import JsonResponse
from django.views.generic import View

# Create your views here.
class AuthenticateView(View):
    """Authenticate a user."""

    def post(request):
        username = request.POST['username']
        password = request.POST['password']
        user = authenticate(username=username, password=password)
        if user is not None:
            if user.is_active:
                # login(request, user)
                # Redirect to a success page.
                return JsonResponse({'message': 'You have successfully authenticated',
                                     'error': False})
            else:
                # Return a 'disabled account' error message
               return JsonResponse({'message': 'Your account is not active',
                                    'error': True})
        else:
            # Return an 'invalid login' error message.
            return JsonResponse({'message': 'Credentials are invalid',
                                 'error': True})

