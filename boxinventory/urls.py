from django.urls import path

from . import views

app_name = 'boxinventory'
urlpatterns = [
    path('', views.index, name='index'),
    path('qr-sheet', views.qr_sheet, name='qr-sheet'),
    path('<int:box_id>', views.box, name='box'),
    path('<int:box_id>/qr', views.qr, name='qr'),
    path('<int:box_id>/add', views.addplant, name='addplant'),
    path('<int:box_id>/set-qr-applied', views.set_qr_applied, name='set-qr-applied'),
]
