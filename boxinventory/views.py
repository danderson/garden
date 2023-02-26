from django.shortcuts import render
from django.http import HttpResponse
from .models import Box

def index(request):
    boxes = Box.objects.order_by('name')
    ctx = {'boxes': boxes}
    return render(request, 'boxinventory/index.html', ctx)
