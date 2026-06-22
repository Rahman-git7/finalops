# Guide de déploiement FinalOps — Homelab Kubernetes

## Prérequis matériel

| Ressource | Minimum | Recommandé |
|-----------|---------|------------|
| RAM | 2 Go | 4 Go |
| CPU | 2 cœurs | 4 cœurs |
| Disque | 20 Go | 50 Go |
| OS | Ubuntu 22.04 LTS | Ubuntu 22.04 LTS |

---

## Étape 1 — Configurer l'inventaire Ansible

Édite `infrastructure/ansible/inventory/homelab.ini` :
```ini
[homelab]
homelab-server ansible_host=192.168.1.XXX ansible_user=ubuntu
```

Remplace `192.168.1.XXX` par l'IP de ton serveur.

---

## Étape 2 — Installer k3s

```bash
ansible-playbook -i infrastructure/ansible/inventory/homelab.ini \
  infrastructure/ansible/playbooks/01-install-k3s.yml
```

Ensuite, configure kubectl sur ta machine locale :
```bash
cp ~/.kube/finalops-homelab.yaml ~/.kube/config
# Édite le fichier et remplace '127.0.0.1' par l'IP du serveur
export KUBECONFIG=~/.kube/finalops-homelab.yaml
kubectl get nodes  # doit afficher le node Ready
```

---

## Étape 3 — Créer les secrets (IMPORTANT)

Sur ta machine locale avec kubectl configuré :

```bash
kubectl create secret generic finalops-secrets \
  --namespace finalops \
  --from-literal=JWT_SECRET=$(openssl rand -hex 32) \
  --from-literal=POSTGRES_PASSWORD=$(openssl rand -hex 16) \
  --from-literal=AUTH_DB_PASSWORD=$(openssl rand -hex 16) \
  --from-literal=EXERCISE_DB_PASSWORD=$(openssl rand -hex 16) \
  --from-literal=WORKOUT_DB_PASSWORD=$(openssl rand -hex 16)
```

> Ces secrets ne sont JAMAIS dans git. Tu les recréés si tu rebuild le cluster.

---

## Étape 4 — Installer ArgoCD

```bash
ansible-playbook -i infrastructure/ansible/inventory/homelab.ini \
  infrastructure/ansible/playbooks/02-install-argocd.yml
```

Accès à l'UI ArgoCD depuis ta machine locale :
```bash
kubectl port-forward svc/argocd-server -n argocd 8080:443
# Ouvre https://localhost:8080
# Login: admin / password affiché par le playbook
```

---

## Étape 5 — Déployer FinalOps via ArgoCD

```bash
kubectl apply -f argocd/app-of-apps.yaml
```

ArgoCD va automatiquement syncer tous les manifests depuis GitHub et déployer l'app.

Surveille le déploiement :
```bash
kubectl get pods -n finalops -w
```

---

## Étape 6 — Accéder à l'app

Ajoute dans `/etc/hosts` sur ta machine cliente :
```
192.168.1.XXX  finalops.homelab
```

Ouvre **http://finalops.homelab**

---

## Workflow GitOps (une fois déployé)

```
Tu fais un commit sur main
        ↓
GitHub Actions build + push l'image sur GHCR
        ↓
GitHub Actions patch le deployment.yaml avec le nouveau SHA
        ↓
ArgoCD détecte le changement dans git
        ↓
ArgoCD déploie automatiquement la nouvelle version
```

---

## Commandes utiles

```bash
# Voir tous les pods
kubectl get pods -n finalops

# Logs d'un service
kubectl logs -n finalops -l app=auth-service -f

# Voir les events en cas de problème
kubectl get events -n finalops --sort-by='.lastTimestamp'

# Forcer un sync ArgoCD
kubectl patch application finalops-infrastructure -n argocd \
  --type merge -p '{"operation":{"sync":{}}}'

# Accéder à ArgoCD UI
kubectl port-forward svc/argocd-server -n argocd 8080:443
```
