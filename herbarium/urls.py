from django.urls import path

from . import views

app_name = 'herbarium'
urlpatterns = [
    path('', views.index),
    path('now/', views.now),
    path('missing/', views.missingWindows, name='missingWindows'),
    path('missing/<int:id>', views.updateMissing, name='updateMissing'),
]
