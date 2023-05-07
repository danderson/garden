from django import forms
from .models import BoxContent

class AddPlantForm(forms.ModelForm):
    class Meta:
        model = BoxContent
        fields = ['name', 'planted']
