from django import forms
from .models import BoxContent

class AddPlantForm(forms.ModelForm):
    class Meta:
        model = BoxContent
        fields = ['name', 'planted']

class RemovePlantForm(forms.Form):
    id = forms.IntegerField(widget=forms.HiddenInput())
    remove = forms.BooleanField(required=False)

RemovePlantFormset = forms.formset_factory(RemovePlantForm, extra=0)
