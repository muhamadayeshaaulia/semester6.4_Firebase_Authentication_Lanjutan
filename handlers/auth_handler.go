package handlers
import(
	"net/http"
 	"time"
 	"github.com/gin-gonic/gin"
 	"github.com/muhamadayeshaaulia/gin-firebase-backend/services"
)
type AuthHandler struct {
 authService *services.AuthService
}
