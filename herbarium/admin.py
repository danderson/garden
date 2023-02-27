from django.contrib import admin

from .models import Family, Plant, Variety, PlantName

class VarietyInline(admin.StackedInline):
    model = Variety
    extra = 1

class PlantNameInline(admin.StackedInline):
    model = PlantName
    extra = 1

class PlantAdmin(admin.ModelAdmin):
    inlines = [PlantNameInline, VarietyInline]

admin.site.register(Family)
admin.site.register(Plant, PlantAdmin)
admin.site.register(Variety)
