from django.conf.urls import patterns, include, url
from django.contrib import admin

import api.urls
import server.views

urlpatterns = patterns(
    '',
    url(r'^$', server.views.HomeView.as_view(), name='home'),
    url(r'^admin/', include(admin.site.urls)),
    url(r'^api/', include(api.urls, namespace='api')),
)
