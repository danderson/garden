from django.shortcuts import render
from django.http import HttpResponse
from .models import Box

def index(request):
    boxes = Box.objects.all()
    ctx = {'boxes': boxes}
    return render(request, 'boxinventory/index.html', ctx)
