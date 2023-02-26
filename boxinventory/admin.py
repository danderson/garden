from django.contrib import admin

from .models import Box, BoxContent

class BoxContentInline(admin.StackedInline):
    model = BoxContent
    extra = 1

class BoxAdmin(admin.ModelAdmin):
    inlines = [BoxContentInline]

admin.site.register(Box, BoxAdmin)
