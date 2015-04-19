from django.conf.urls import patterns, url
# from django.contrib import admin

from api.views import LoginView, RegisterView, WorkoutView

urlpatterns = patterns(
    '',
    url(r'^auth/$', views.LoginView.as_view(), name='login'),
    url(r'^register/$', views.RegisterView.as_view(), name='register'),
    url(r'^workout/$', views.WorkoutView.as_view(), name='workout'),
)
