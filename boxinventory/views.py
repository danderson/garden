from django.shortcuts import render
from django.http import HttpResponse
from .models import Box
import qrcode
from django.urls import reverse

def index(request):
    boxes = Box.objects.order_by('name')
    ctx = {'boxes': boxes}
    return render(request, 'boxinventory/index.html', ctx)

def box(request, box_id):
    box = Box.objects.get(pk=box_id)
    ctx = {'box': box}
    return render(request, 'boxinventory/box.html', ctx)

def qr(request, box_id):
    code = qrcode.QRCode(
        error_correction=qrcode.constants.ERROR_CORRECT_H,
        box_size=10,
    )
    code.add_data(reverse('boxinventory:box', args=(box_id,)))
    code.make(fit=True)
    img = code.make_image(fill_color="black", back_color="white")
    resp = HttpResponse(content_type="image/png")
    img.save(resp, "png")
    return resp
