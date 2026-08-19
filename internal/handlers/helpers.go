package handlers

import "github.com/gin-gonic/gin"

func errJSON(c *gin.Context, s int, e error) { c.AbortWithStatusJSON(s, gin.H{"error": e.Error()}) }
