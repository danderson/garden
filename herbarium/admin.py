from django.contrib import admin

from .models import Family, Plant, Variety

class VarietyInline(admin.StackedInline):
    model = Variety
    extra = 1

class PlantAdmin(admin.ModelAdmin):
    inlines = [VarietyInline]

admin.site.register(Family)
admin.site.register(Plant, PlantAdmin)
admin.site.register(Variety)
