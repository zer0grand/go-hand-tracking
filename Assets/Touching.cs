using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class Touching : MonoBehaviour {

  public TextMesh textEnter;

  private Color color1;
  private Color color2;

  void Start() {
    color1 = new Color(255,255,0);
    color2 = new Color(0,255,255);
  }

  void OnCollisionEnter(Collision col) {
    if (col.gameObject.tag == "key") {
      col.gameObject.GetComponent<Renderer>().material.color = color1;
      string name = col.gameObject.name;
      if (name.Length == 1) {
        textEnter.text += name;
      } else if (name == "delete") {
        if (textEnter.text.Length > 0) {
          textEnter.text = textEnter.text.Substring(0, textEnter.text.Length-1);
        }
      }
    }
  }

  void OnCollisionExit(Collision col) {
    if (col.gameObject.tag == "key") {
      col.gameObject.GetComponent<Renderer>().material.color = color2;
    }
  }
}
