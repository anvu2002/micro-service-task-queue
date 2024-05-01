from dataclasses import dataclass

@dataclass
class ImageScore:
    score: float
    url: str
    description: str

class TTSResponse:
    
    def __init__(self, status, msg):
        self.status = status
        self.msg = msg

    def to_dict(self):
        return {
            "status": self.status,
            "msg": self.msg
        }