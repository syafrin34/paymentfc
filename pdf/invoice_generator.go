package pdf

import (
	"fmt"
	"paymentfc/models"

	"github.com/phpdave11/gofpdf"
)

func GenerateInvoicePDF(payment models.Payment, outputPath string) error {
	pdf := gofpdf.New("P", "mm", "A4", "arial")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "[FC] Invoice Details")
	pdf.Ln(20)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Payment ID: #%d", payment.ID))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Order ID: #%d", payment.OrderID))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Total Amount: Rp%.2f", payment.Amount))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Status: %s", payment.Status))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("External ID: %s", payment.ExternalID))
	return pdf.OutputFileAndClose(outputPath)
}
