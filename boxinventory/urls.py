from django.urls import path

from . import views

app_name = 'boxinventory'
urlpatterns = [
    path('', views.index, name='index'),
    path('box/<int:box_id>', views.box, name='box'),
    path('box/<int:box_id>/qr', views.qr, name='qr'),
]
