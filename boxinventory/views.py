from django.shortcuts import render, redirect
from django.http import HttpResponse
from .models import Box, BoxContent
from .forms import AddPlantForm, RemovePlantFormset
import qrcode
from django.urls import reverse
import datetime

def index(request):
    boxes = Box.objects.order_by('name')
    ctx = {'boxes': boxes}
    return render(request, 'boxinventory/index.html', ctx)

def box(request, box_id):
    box = Box.objects.get(id=box_id)
    contents = BoxContent.objects.filter(box__id=box_id, removed__isnull=True).order_by('name', '-planted')
    removed = BoxContent.objects.filter(box__id=box_id, removed__isnull=False).order_by('-planted', 'name')
    ctx = {
        'box': box,
        'contents': contents,
        'removed': removed,
    }
    return render(request, 'boxinventory/box.html', ctx)

def addplant(request, box_id):
    if request.method == 'POST':
        form = AddPlantForm(request.POST)
        if form.is_valid():
            box = Box.objects.get(pk=box_id)
            ct = form.save(commit=False)
            ct.box = box
            ct.save()
            return redirect('boxinventory:box', box_id)
    else:
        form = AddPlantForm()

    return render(request, 'boxinventory/addplant.html', {'form': form})

def remove_plants(request, box_id):
    if request.method == 'POST':
        form = RemovePlantFormset(request.POST)
        if form.is_valid():
            for obj in form.cleaned_data:
                if obj['remove']:
                    bc = BoxContent.objects.get(id=obj['id'])
                    bc.removed = datetime.date.today()
                    bc.save()
            return redirect('boxinventory:box', box_id)
    else:
        content = BoxContent.objects.filter(box__id=box_id, removed=None).order_by('-planted', 'name')
        initial = [{'id': c.id, 'name': c.name, 'remove': False} for c in content]
        form = RemovePlantFormset(initial=initial)
        for f, c in zip(form.forms, content):
            f.fields['remove'].label = c.dated_name

    return render(request, 'boxinventory/removeplants.html', {'form': form})

def set_qr_applied(request, box_id):
    box = Box.objects.get(pk=box_id)
    box.qr_applied = True
    box.save()
    return redirect('boxinventory:box', box_id)

def qr(request, box_id):
    code = qrcode.QRCode(
        error_correction=qrcode.constants.ERROR_CORRECT_H,
        box_size=5,
    )
    code.add_data(request.build_absolute_uri(reverse('boxinventory:box', args=(box_id,))))
    code.make(fit=True)
    img = code.make_image(fill_color="black", back_color="white")
    resp = HttpResponse(content_type="image/png")
    img.save(resp, "png")
    return resp

def qr_sheet(request):
    boxes = Box.objects.filter(want_qr=True, qr_applied=False).order_by('name')
    ctx = {'boxes': boxes}
    return render(request, 'boxinventory/qr-sheet.html', ctx)
