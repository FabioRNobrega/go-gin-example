FROM golang:1.23

WORKDIR /app

# Instala git e outros utilitários
RUN apt-get update && apt-get install -y git vim curl

# Configura usuário Git diretamente
RUN git config --global user.name "FabioRNobrega" && \
    git config --global user.email "fabio.r.nobrega@gmail.com"

CMD ["bash"]
