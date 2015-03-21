"""
views.py
=======
The API & UI views share a common set of core CRUD classes.
The API views implement a JsonResponseMixin.
The UI views implement a SingleObjectTemplateResponseMixin.
"""
# from django.shortcuts import render
# from django.contrib.auth.models import User
from django.contrib.auth import authenticate, login
from django.http import JsonResponse
from django.views.generic import View
from django.views.generic.detail import SingleObjectTemplateResponseMixin


class JSONResponseMixin:

    def render_to_response(self, context, **kwargs):
        """Return a JSONResponse."""
        # TODO Compile and return a JsonResponse
        return JsonResponse(context, **kwargs)

    def form_response(self, form):
        """Return the results of the request."""
        self.object = form.save()
        return self.render_to_response(self.get_context_data(form=form))

    # Overwrite form_*valid to always return a JSONResponse
    form_valid = form_invalid = form_response


class LoginView(View):
    """Authenticate a user."""

    def post(self, request, *args, **kwags):
        username = request.POST.get('username')
        password = request.POST.get('password')
        user = authenticate(username=username, password=password)
        if user is not None:
            if user.is_active:
                login(request, user)
                # Redirect to a success page.
                return JsonResponse({'message':
                                     'You have successfully authenticated',
                                     'error': False})
            else:
                # Return a 'disabled account' error message
                return JsonResponse({'message': 'Your account is not active',
                                    'error': True})
        else:
            # Return an 'invalid login' error message.
            return JsonResponse({'message': 'Credentials are invalid',
                                 'error': True})


class RegisterView(View):
    """Register a new user."""


class WorkoutView(View):
    """Create a new workout."""
