from django.conf.urls import patterns, include, url
# from django.contrib import admin

from . import views

urlpatterns = patterns('',
    url(r'^auth/$', views.AuthenticateView.as_view(), name='auth'),
)
