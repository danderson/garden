from django.urls import path

from . import views

app_name = 'herbarium'
urlpatterns = [
    path('', views.index, name='index'),
    path('now/', views.now, name='now'),
    path('missing/', views.missingWindows, name='missingWindows'),
    path('missing/<int:id>', views.updateMissing, name='updateMissing'),
]
