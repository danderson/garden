import datetime

from django.shortcuts import render, redirect, get_object_or_404
from django.db import transaction
from django import forms

from .models import Plant, PlantingWindow

def index(request):
    return render(request, 'herbarium/index.html')

class WindowForm(forms.Form):
    window_def = forms.CharField(label='Planting window', max_length=200, widget=forms.Textarea)

def missingWindows(request):
    plant = Plant.objects.filter(plantingwindow=None).first()
    return updateMissing(request, plant.id)

@transaction.atomic
def updateMissing(request, id):
    plant = get_object_or_404(Plant, id=id)
    if request.method == 'POST':
        form = WindowForm(request.POST)
        if form.is_valid():
            plant.set_windows_from_str(form.cleaned_data['window_def'])
            return redirect('herbarium:missingWindows')
    else:
        form = WindowForm()
    return render(request, 'herbarium/missingWindows.html', {
        'plant': plant,
        'form': form,
    })
