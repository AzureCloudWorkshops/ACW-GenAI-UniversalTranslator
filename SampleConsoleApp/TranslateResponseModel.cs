namespace ACW_GenAI_UniversalTranslator
{
    public class TranslateResponseModel
    {
        public DetectedLanguage? detectedLanguage { get; set; }
        public Translations[]? translations { get; set; }
    }

    public class DetectedLanguage
    {
        public string? language { get; set; }
        public double? score { get; set; }
    }

    public class Translations
    {
        public string? text { get; set; }
        public string? to { get; set; }
    }
}




